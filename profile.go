package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/model"
	"net/http"
)

func GetProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	sid := mux.Vars(req)["user_id"]
	if len(sid) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	u, err := authCtx.Users.Find(sid)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.NewEncoder(rw).Encode(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func UpdateProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	sid := mux.Vars(req)["user_id"]
	if len(sid) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	u := &model.User{}
	err := json.NewDecoder(req.Body).Decode(u)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Body.Close()

	// don't allow edit ConfirmCodes in this handler
	u.ConfirmCodes = nil
	// check for special privilege
	if u.Privilege != nil || u.BriefGroups != nil || u.Approved != nil {
		_, err := authCtx.ValidCurrentUser(false, nil, []string{"manage_user"})
		if err != nil {
			return http.StatusForbidden, err
		}
	}

	err = authCtx.Users.UpdateDetail(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func ListProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	return http.StatusOK, nil
}
