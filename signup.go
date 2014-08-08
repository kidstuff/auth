package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

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
		return http.StatusBadRequest, errors.New("kidstuff/auth: Pwd and PwdRepeat doesn't match")
	}

	app := true
	if val, err := authCtx.Settings.
		Get("kidstuff.auth.regis.approve_new_user"); err != nil || val != "true" {
		app = false
	}

	u, err := authCtx.Users.Add(credential.Email, credential.PwdRepeat, app)
	if err != nil {
		return http.StatusPreconditionFailed, err
	}

	err = json.NewEncoder(rw).Encode(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if app {
		return http.StatusOK, nil
	}

	mailOpts, err := authCtx.Settings.Get("kidstuff.auth.regis.send_mail_confirm")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var message string

	switch mailOpts {
	case "notice":
		message, err = authCtx.Settings.Get("kidstuff.auth.regis.notice_message")
	case "confirm":
		message, err = authCtx.Settings.Get("kidstuff.auth.regis.confirm_message")
	}

	println(message)

	if err != nil {
		return http.StatusOK, err
	}

	return http.StatusOK, nil
}
