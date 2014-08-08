package auth

import (
	"net/http"
)

func UpdatePassword(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return http.StatusOK, nil
}
