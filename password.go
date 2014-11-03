package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

func OverridePassword(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return updatePassword(authCtx, req, true)
}

func ChangePassword(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return updatePassword(authCtx, req, false)
}

func updatePassword(authCtx *AuthContext, req *http.Request, adminDoingThis bool) (int, error) {
	pwd := struct{ OldPwd, NewPwd, NewPwdRepeat string }{}

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
		err = authCtx.Auth.ComparePassword(pwd.OldPwd, u.Pwd)
		if err != nil {
			return http.StatusPreconditionFailed, ErrInvalidCredential
		}
	}

	err = authCtx.Auth.UpdateUserDetail(*u.Id, &pwd.NewPwd, nil, nil, nil, nil, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
