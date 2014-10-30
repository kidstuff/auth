package auth

import (
	"encoding/json"
	"github.com/kidstuff/auth/authmodel"
	"net/http"
	"strconv"
	"strings"
)

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
