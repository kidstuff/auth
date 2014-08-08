package auth

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
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

	// TODO(!): think about send cofirm email
	return http.StatusOK, nil
}

func Activate(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	sid := vars["user_id"]

	code := req.FormValue("code")
	if len(sid) == 0 || len(code) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	u, err := authCtx.Users.Find(sid)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if ok := u.ValidConfirmCode("activate", code, false, true); !ok {
		return http.StatusPreconditionFailed, ErrInvalidActiveCode
	}

	err = authCtx.Users.UpdateDetail(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	rw.Write([]byte(`{"Message":"Account activated"}`))
	return http.StatusOK, nil
}
