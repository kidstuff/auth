package auth

import (
	"net/http"
)

func GetProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return http.StatusOK, nil
}

func UpdateProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return http.StatusOK, nil
}

func ListProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return http.StatusOK, nil
}
