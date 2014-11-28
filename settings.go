package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

// UpdateSettings document: http://kidstuff.github.io/swagger/#!/default/settings_patch
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

// GetSettings document: http://kidstuff.github.io/swagger/#!/default/settings_get
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

// DeleteSettings document: http://kidstuff.github.io/swagger/#!/default/settings_delete
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
