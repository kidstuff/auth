package auth

import (
	"encoding/json"
	"errors"
	"github.com/kidstuff/auth/model"
	"net/http"
	"time"
)

func GetToken(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	grantType := req.FormValue("grant_type")
	email := req.FormValue("email")
	password := req.FormValue("password")

	// TODO: more detail error message
	if len(grantType) == 0 || len(email) == 0 || len(password) == 0 {
		return http.StatusBadRequest, errors.New("kidstuff/auth: grant_type, email and password need to be set.")
	}

	if grantType != "password" {
		return http.StatusBadRequest, errors.New("kidstuff/auth: Only support grant_type=password")
	}

	user, err := authCtx.Users.FindByEmail(email)
	if err != nil {
		return http.StatusUnauthorized, ErrInvalidCredential
	}

	err = user.ComparePassword(password)
	if err != nil {
		return http.StatusUnauthorized, ErrInvalidCredential
	}

	token, err := authCtx.Users.Login(user.Id, OnlineThreshold)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	inf := struct {
		User        *model.User
		ExpiredOn   time.Time
		AccessToken string
	}{user, time.Now().Add(OnlineThreshold), token}

	return http.StatusOK, json.NewEncoder(rw).Encode(&inf)
}
