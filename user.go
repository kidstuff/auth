package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/authmodel"
	"net/http"
	"strconv"
	"strings"
)

func findUser(ctx *AuthContext, req *http.Request) (*authmodel.User, int, error) {
	sid := mux.Vars(req)["user_id"]
	if len(sid) == 0 {
		return nil, http.StatusBadRequest, ErrInvalidId
	}

	u, err := ctx.Auth.FindUser(sid)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return u, http.StatusOK, nil
}

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

func DeleteUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	u, stt, err := findUser(authCtx, req)
	if err != nil {
		return stt, err
	}

	err = authCtx.Auth.DeleteUser(*u.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

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

func RemoveGroupFromUser(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	params := mux.Vars(req)
	sid := params["user_id"]
	gid := params["group_id"]
	if len(sid) == 0 || len(gid) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	u, err := authCtx.Auth.FindUser(sid)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	ids := make([]string, 0, len(u.Groups))
	for _, group := range u.Groups {
		if *group.Id != gid {
			ids = append(ids, *group.Id)
		}
	}

	err = authCtx.Auth.UpdateUserDetail(*u.Id, nil, nil, nil, nil, nil, ids)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
