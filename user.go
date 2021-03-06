package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/authmodel"
	"net/http"
	"strconv"
	"strings"
)

// CreateUser handle create new user accoutn action. Require "manage_user" privilege.
// Details: http://kidstuff.github.io/swagger/#!/default/users_post
func CreateUser(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	info := struct {
		Email     string
		Pwd       string
		PwdRepeat string
		Approved  bool
	}{}

	err := json.NewDecoder(req.Body).Decode(&info)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if info.Pwd != info.PwdRepeat {
		return http.StatusBadRequest, ErrPwdMismatch
	}

	u, err := ctx.Auth.AddUser(info.Email, info.PwdRepeat, info.Approved)
	if err != nil {
		return http.StatusPreconditionFailed, err
	}

	err = json.NewEncoder(rw).Encode(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func findUser(ctx *AuthContext, req *http.Request) (*authmodel.User, int, error) {
	sid := mux.Vars(req)["user_id"]
	if len(sid) == 0 {
		return nil, http.StatusBadRequest, ErrInvalidId
	}

	u, err := ctx.Auth.FindUser(sid)
	if err != nil {
		if err == authmodel.ErrNotFound {
			return nil, http.StatusNotFound, err
		}

		return nil, http.StatusInternalServerError, err
	}

	return u, http.StatusOK, nil
}

// GetUser return full user account data. Require "manage_user" privilege.
// Details: http://kidstuff.github.io/swagger/#!/default/users_user_id_get
func GetUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	u, stt, err := findUser(authCtx, req)
	if err != nil {
		return stt, err
	}

	err = json.NewEncoder(rw).Encode(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteUser handle delete user account action. Require "manage_user" privilege.
// Details: http://kidstuff.github.io/swagger/#!/default/users_user_id_delete
func DeleteUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	sid := mux.Vars(req)["user_id"]
	if len(sid) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	err := authCtx.Auth.DeleteUser(sid)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// UpdateUserProfile handle user's profile update action. Require "owner" or "manage_user" privilege.
// Details: http://kidstuff.github.io/swagger/#!/default/users_user_id_profile_patch
func UpdateUserProfile(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	u, stt, err := findUser(authCtx, req)
	if err != nil {
		return stt, err
	}

	p := &authmodel.Profile{}
	err = json.NewDecoder(req.Body).Decode(p)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Body.Close()

	err = authCtx.Auth.UpdateUserDetail(*u.Id, nil, nil, nil, nil, p, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// ListUser handle user list action. Require "manage_user" privilege.
// Details: http://kidstuff.github.io/swagger/#!/default/users_get
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

	var groupIds []string
	if groups := req.FormValue("groups"); len(groups) > 0 {
		groupIds = strings.Split(groups, ",")
	}

	users, err := authCtx.Auth.FindAllUser(limit, offsetId, selectFields, groupIds)
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
		Users []*authmodel.User
		Next  string
	}{users, next}

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// UpdateApprovedStatus handle user arppoval status update action. Require "manage_user" privilege.
// Details: http://kidstuff.github.io/swagger/#!/default/users_user_id_approve_put
func UpdateApprovedStatus(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	u, stt, err := findUser(authCtx, req)
	if err != nil {
		return stt, err
	}

	app := struct{ Approved bool }{}
	err = json.NewDecoder(req.Body).Decode(&app)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Body.Close()

	err = authCtx.Auth.UpdateUserDetail(*u.Id, nil, &app.Approved, nil, nil, nil, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// AddGroupToUser add group to user. Require "manage_user" privilege.
// Details: http://kidstuff.github.io/swagger/#!/default/users_user_id_groups_put
func AddGroupToUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	u, stt, err := findUser(authCtx, req)
	if err != nil {
		return stt, err
	}

	g := struct{ Id string }{}
	err = json.NewDecoder(req.Body).Decode(&g)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Body.Close()

	for _, group := range u.Groups {
		if *group.Id == g.Id {
			return http.StatusOK, nil
		}
	}

	ids := make([]string, 1, len(u.Groups)+1)
	ids[0] = g.Id
	for _, group := range u.Groups {
		ids = append(ids, *group.Id)
	}

	err = authCtx.Auth.UpdateUserDetail(*u.Id, nil, nil, nil, nil, nil, ids)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemoveGroupFromUser document: http://kidstuff.github.io/swagger/#!/default/users_user_id_groups_group_id_delete
func RemoveGroupFromUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	params := mux.Vars(req)
	sid := params["user_id"]
	gid := params["group_id"]
	if len(sid) == 0 || len(gid) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	u, err := authCtx.Auth.FindUser(sid)
	if err != nil {
		if err == authmodel.ErrNotFound {
			return http.StatusNotFound, err
		}

		return http.StatusInternalServerError, err
	}

	n := len(u.Groups)
	ids := make([]string, 0, n)
	for _, group := range u.Groups {
		if *group.Id != gid {
			ids = append(ids, *group.Id)
		}
	}

	// if the only group was removed, then ids must be a zeroed slice (not an nil slice).
	// That allow UpdateUserDetail take the action
	if len(ids) == 0 && n == 1 {
		ids = []string{}
	}

	err = authCtx.Auth.UpdateUserDetail(*u.Id, nil, nil, nil, nil, nil, ids)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
