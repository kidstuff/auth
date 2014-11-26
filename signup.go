package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// TODO(!): figure out how to properly cancel the email send when time out

func sendWelcomeMail(ctx *AuthContext, email string) error {
	if val, err := ctx.Settings.Get("auth_send_welcome_email"); val != "true" {
		return nil
	} else if err != nil {
		return err
	}

	mailSettings, err := ctx.Settings.GetMulti([]string{
		"auth_full_path",
		"auth_welcome_email_subject",
		"auth_welcome_email_message",
		"auth_email_from",
	})
	if err != nil {
		return err
	}

	return ctx.Notifications.SendMail(ctx, mailSettings["auth_welcome_email_subject"],
		mailSettings["auth_welcome_email_message"],
		mailSettings["auth_email_from"], email)
}

func sendActivateMail(ctx *AuthContext, id, email, code string) error {
	if val, err := ctx.Settings.Get("auth_send_activate_email"); val != "true" {
		return nil
	} else if err != nil {
		return err
	}

	mailSettings, err := ctx.Settings.GetMulti([]string{
		"auth_full_path",
		"auth_activate_page",
		"auth_activate_email_subject",
		"auth_activate_email_message",
		"auth_email_from",
	})
	if err != nil {
		return err
	}

	activeURL := fmt.Sprintf("%s/users/%s/activate?code=%s", mailSettings["auth_full_path"], id, code)
	return ctx.Notifications.SendMail(ctx, mailSettings["auth_activate_email_subject"],
		fmt.Sprintf(mailSettings["auth_activate_email_message"], activeURL),
		mailSettings["auth_email_from"], email)
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

	status := 200
	if app {
		err = sendWelcomeMail(authCtx, *u.Email)
		if err != nil {
			authCtx.Logs.Errorf("Wellcome mail failed: %s", err)
		}
	} else {
		err = sendActivateMail(authCtx, *u.Id, *u.Email, u.ConfirmCodes["activate"])
		if err != nil {
			authCtx.Logs.Errorf("Active mail failed: %s", err)
		}
		status = http.StatusAccepted
	}

	json.NewEncoder(rw).Encode(u)
	return status, nil
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

	err = sendWelcomeMail(authCtx, *u.Email)
	if err != nil {
		authCtx.Logs.Errorf("Wellcome mail failed: %s", err)
	}

	activate_redirect, err := authCtx.Settings.Get("auth_activate_redirect")
	if err != nil {
		authCtx.Logs.Errorf("Error when fetching 'auth_activate_redirect' settings")
		rw.Write([]byte(`{"Message":"Account activated"}`))
	} else {
		http.Redirect(rw, req, activate_redirect, http.StatusSeeOther)
		rw.Write([]byte(`{"Message":"Account activated", "RedirectTo": "` + activate_redirect + `"}`))
	}

	return http.StatusOK, nil
}
