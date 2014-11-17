package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func sendWelcomeMail(authCtx *AuthContext, email string) (int, error) {
	if val, err := authCtx.Settings.
		Get("auth_send_welcome_email"); err != nil || val != "true" {
		return http.StatusOK, nil
	}

	mailSettings, err := authCtx.Settings.GetMulti([]string{
		"auth_full_path",
		"auth_welcome_email_subject",
		"auth_welcome_email_message",
		"auth_email_from",
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = DEFAULT_NOTIFICATOR.SendMail(mailSettings["auth_welcome_email_subject"],
		mailSettings["auth_welcome_email_message"],
		mailSettings["auth_email_from"], email)
	if err != nil {
		authCtx.Logs.Errorf("mail send error:%s", err)
	}

	return http.StatusOK, nil
}

func SignUp(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	credential := struct {
		Email     string
		Pwd       string
		PwdRepeat string
	}{}

	err := json.NewDecoder(req.Body).Decode(&credential)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if credential.Pwd != credential.PwdRepeat {
		return http.StatusBadRequest, ErrPwdMismatch
	}

	app := true
	if val, err := authCtx.Settings.
		Get("auth_approve_new_user"); err != nil || val != "true" {
		app = false
	}

	u, err := authCtx.Auth.AddUser(credential.Email, credential.PwdRepeat, app)
	if err != nil {
		return http.StatusPreconditionFailed, err
	}

	err = json.NewEncoder(rw).Encode(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if app {
		return sendWelcomeMail(authCtx, *u.Email)
	} else {
		if val, err := authCtx.Settings.
			Get("auth_send_activate_email"); err != nil || val != "true" {
			return http.StatusOK, nil
		}

		mailSettings, err := authCtx.Settings.GetMulti([]string{
			"auth_full_path",
			"auth_activate_page",
			"auth_activate_email_subject",
			"auth_activate_email_message",
			"auth_email_from",
		})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		activeURL := fmt.Sprintf(mailSettings["auth_activate_page"], *u.Id, u.ConfirmCodes["activate"])
		err = DEFAULT_NOTIFICATOR.SendMail(mailSettings["auth_activate_email_subject"],
			fmt.Sprintf(mailSettings["auth_activate_email_message"], activeURL),
			mailSettings["auth_email_from"], *u.Email)
		if err != nil {
			authCtx.Logs.Errorf("Send mail failed:%s", err)
		}
	}

	return http.StatusOK, nil
}

func Activate(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	sid := vars["user_id"]

	code := req.FormValue("code")
	if len(sid) == 0 || len(code) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	u, err := authCtx.Auth.FindUser(sid)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if ok := u.ValidConfirmCode("activate", code, false, true); !ok {
		return http.StatusPreconditionFailed, ErrInvalidActiveCode
	}

	t := true
	err = authCtx.Auth.UpdateUserDetail(*u.Id, nil, &t, nil, nil, nil, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	stt, err := sendWelcomeMail(authCtx, *u.Email)
	if err != nil {
		return stt, err
	}

	rw.Write([]byte(`{"Message":"Account activated"}`))
	return http.StatusOK, nil
}
