package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/authmodel"
	"net/http"
	"strconv"
	"strings"
)

// CreateGroup document: http://kidstuff.github.io/swagger/#!/default/groups_post
func CreateGroup(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	groups := &authmodel.Group{}
	err := json.NewDecoder(req.Body).Decode(groups)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Body.Close()

	groups, err = authCtx.Auth.AddGroupDetail(*groups.Name, groups.Privileges, groups.Info)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.NewEncoder(rw).Encode(groups)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func findGroup(ctx *AuthContext, req *http.Request) (*authmodel.Group, int, error) {
	sid := mux.Vars(req)["group_id"]
	if len(sid) == 0 {
		return nil, http.StatusBadRequest, ErrInvalidId
	}

	g, err := ctx.Auth.FindGroup(sid)
	if err != nil {
		if err == authmodel.ErrNotFound {
			return nil, http.StatusNotFound, err
		}

		return nil, http.StatusInternalServerError, err
	}

	return g, http.StatusOK, nil
}

// UpdateGroup document: http://kidstuff.github.io/swagger/#!/default/groups_group_id_patch
func UpdateGroup(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	g, stt, err := findGroup(ctx, req)
	if err != nil {
		return stt, err
	}

	g2 := authmodel.Group{}
	err = json.NewDecoder(req.Body).Decode(&g2)
	req.Body.Close()
	if err != nil {
		return http.StatusBadRequest, err
	}

	ctx.Logs.Debugf("%#v", g2)

	err = ctx.Auth.UpdateGroupDetail(*g.Id, g2.Privileges, g2.Info)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetGroup document: http://kidstuff.github.io/swagger/#!/default/groups_group_id_get
func GetGroup(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	g, stt, err := findGroup(ctx, req)
	if err != nil {
		return stt, err
	}

	err = json.NewEncoder(rw).Encode(g)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteGroup document: http://kidstuff.github.io/swagger/#!/default/groups_group_id_delete
func DeleteGroup(ctx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	sid := mux.Vars(req)["group_id"]
	if len(sid) == 0 {
		return http.StatusBadRequest, ErrInvalidId
	}

	err := ctx.Auth.DeleteGroup(sid)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// ListGroup document: http://kidstuff.github.io/swagger/#!/default/groups_get
func ListGroup(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request) (int, error) {
	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = -1
	}

	offsetId := req.FormValue("offset")
	var selectFields []string
	if slt := req.FormValue("select"); len(slt) > 0 {
		selectFields = strings.Split(slt, ",")
	}

	groups, err := authCtx.Auth.FindAllGroup(limit, offsetId, selectFields)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	next := req.URL.String()
	if err != nil {
		q := req.URL.Query()
		q.Set("offset", *groups[len(groups)-1].Id)
		req.URL.RawQuery = q.Encode()
		next = req.URL.String()
	}

	response := struct {
		Groups []*authmodel.Group
		Next   string
	}{groups, next}

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
