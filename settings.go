package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

func UpdateSettings(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	settings := map[string]string{}
	err := json.NewDecoder(req.Body).Decode(&settings)
	req.Body.Close()
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = ctx.Settings.SetMulti(settings)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetSettings(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	keys := strings.Split(req.FormValue("keys"), ",")
	if len(keys) == 0 {
		return http.StatusBadRequest, ErrNoKeyProvided
	}

	settings, err := ctx.Settings.GetMulti(keys)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.NewEncoder(rw).Encode(settings)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func DeleteSettings(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	keys := strings.Split(req.FormValue("keys"), ",")
	if len(keys) == 0 {
		return http.StatusBadRequest, ErrNoKeyProvided
	}

	err := ctx.Settings.UnSetMulti(keys)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
