package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/kidstuff/auth/authmodel"
	"net/http"
	"strings"
)

// OverridePassword document: http://kidstuff.github.io/swagger/#!/default/users_user_id_password_override_put
func OverridePassword(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return updatePassword(authCtx, req, true, false)
}

// ChangePassword document: http://kidstuff.github.io/swagger/#!/default/users_user_id_password_put
func ChangePassword(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return updatePassword(authCtx, req, false, false)
}

// ResetPassword document: http://kidstuff.github.io/swagger/#!/default/users_user_id_password_reset_put
func ResetPassword(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return updatePassword(authCtx, req, false, true)
}

func updatePassword(authCtx *AuthContext, req *http.Request, adminDoingThis, isResetting bool) (int, error) {
	pwd := struct{ OldPwd, NewPwd, NewPwdRepeat, ResetCode string }{}

	err := json.NewDecoder(req.Body).Decode(&pwd)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Body.Close()

	if pwd.NewPwd != pwd.NewPwdRepeat {
		return http.StatusBadRequest, errors.New("kidstuff/auth: passsword mismatch")
	}

	u, stt, err := findUser(authCtx, req)
	if err != nil {
		return stt, err
	}

	if !adminDoingThis {
		if isResetting {
			if len(pwd.ResetCode) == 0 {
				return http.StatusPreconditionFailed, ErrInvalidResetCode
			}

			if !u.ValidConfirmCode("password_reset", pwd.ResetCode, false, true) {
				return http.StatusPreconditionFailed, ErrInvalidResetCode
			}
		} else {
			if len(pwd.OldPwd) == 0 {
				return http.StatusPreconditionFailed, ErrInvalidCredential
			}

			err := authCtx.Auth.ComparePassword(pwd.OldPwd, u.Pwd)
			if err != nil {
				return http.StatusPreconditionFailed, ErrInvalidCredential
			}
		}
	}

	err = authCtx.Auth.UpdateUserDetail(*u.Id, &pwd.NewPwd, nil, nil, u.ConfirmCodes, nil, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func sendResetMail(ctx *AuthContext, id, email, code string) error {
	mailSettings, err := ctx.Settings.GetMulti([]string{
		"auth_full_path",
		"auth_reset_email_subject",
		"auth_reset_email_message",
		"auth_email_from",
	})
	if err != nil {
		return err
	}

	activeURL := fmt.Sprintf("%s/password/reset/%s?code=%s", mailSettings["auth_full_path"], id, code)
	return ctx.Notifications.SendMail(ctx, mailSettings["auth_reset_email_subject"],
		fmt.Sprintf(mailSettings["auth_reset_email_message"], activeURL),
		mailSettings["auth_email_from"], email)
}

func CreatePasswordResetIssue(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	inf := struct {
		Resend bool
		Email  string
	}{}

	err := json.NewDecoder(req.Body).Decode(&inf)
	req.Body.Close()
	if err != nil || len(inf.Email) == 0 {
		return http.StatusBadRequest, err
	}

	u, err := ctx.Auth.FindUserByEmail(inf.Email)
	if err == authmodel.ErrNotFound {
		return http.StatusPreconditionFailed, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(u.ConfirmCodes["password_reset"]) > 0 && !inf.Resend {
		return http.StatusNotAcceptable, errors.New("kidstuff/auth: password reset has been request")
	}

	if u.ConfirmCodes == nil {
		u.ConfirmCodes = map[string]string{}
	}

	u.ConfirmCodes["password_reset"] = strings.Trim(base64.URLEncoding.
		EncodeToString(securecookie.GenerateRandomKey(64)), "=")

	err = ctx.Auth.UpdateUserDetail(*u.Id, nil, nil, nil, u.ConfirmCodes, nil, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = sendResetMail(ctx, *u.Id, *u.Email, u.ConfirmCodes["password_reset"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func RedirectPasswordReset(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	sid := mux.Vars(req)["user_id"]
	if len(sid) == 0 || len(req.FormValue("code")) == 0 {
		return http.StatusBadRequest, ErrInvalidCredential
	}

	mailSettings, err := ctx.Settings.GetMulti([]string{
		"auth_reset_redirect",
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	urlStr := fmt.Sprintf(mailSettings["auth_reset_redirect"], sid, req.FormValue("code"))
	http.Redirect(rw, req, urlStr, http.StatusSeeOther)
	return http.StatusSeeOther, nil
}
