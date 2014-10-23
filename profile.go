package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/model"
	"net/http"
	"strconv"
	"strings"
)

func GetUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
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

func UpdateUserProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	sid := mux.Vars(req)["user_id"]
	if len(sid) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	u, err := authCtx.Users.Find(sid)
	if err != nil {
		return http.StatusNotFound, err
	}

	p := &model.Profile{}
	err = json.NewDecoder(req.Body).Decode(p)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Body.Close()

	err = authCtx.Users.UpdateDetail(*u.Id, nil, nil, nil, nil, p, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func ListUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = -1
	}

	offsetId := req.FormValue("offset")
	var selectFields []string
	if slt := req.FormValue("select"); len(slt) > 0 {
		selectFields = strings.Split(slt, ",")
	}

	users, err := authCtx.Users.FindAll(limit, offsetId, selectFields)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	next := req.URL.String()
	if err != nil {
		q := req.URL.Query()
		q.Set("offset", *users[len(users)-1].Id)
		req.URL.RawQuery = q.Encode()
		next = req.URL.String()
	}

	response := struct {
		Users []*model.User
		Next  string
	}{users, next}

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
