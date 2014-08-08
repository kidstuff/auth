package auth

import (
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/model"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	OnlineThreshold = time.Hour
)

type ctxKey int

const (
	userTokenKey ctxKey = iota
	userCurrentKey
)

func Serve(router *mux.Router) {
	if HandlerRegister == nil || EqualIdChecker == nil {
		panic("kidstuff/auth: need to be initialed by a mngr")
	}

	router.Handle("/tokens", HandlerRegister(GetToken, false, nil, nil))
	router.Handle("/users/{user_id}/activate", HandlerRegister(Activate, false, nil, nil))
	router.Handle("/users/{user_id}/password", HandlerRegister(UpdatePassword, true, nil, nil)).Methods("PUT")
	router.Handle("/users/{user_id}", HandlerRegister(GetProfile, false, nil, nil)).Methods("GET")
	router.Handle("/users/{user_id}", HandlerRegister(UpdateProfile, true, nil, nil)).Methods("PATCH")
	router.Handle("/users", HandlerRegister(ListProfile, false, nil, nil))

}

type AuthContext struct {
	Users    model.UserManager
	Groups   model.GroupManager
	Settings model.Configurator
}

type HandleFunc func(*AuthContext, http.ResponseWriter, *http.Request) (int, error)

var (
	HandlerRegister func(fn HandleFunc, owner bool, groups, pri []string) http.Handler
	EqualIdChecker  func(interface{}, string) bool
)

type BasicMngrHandler struct {
	AuthContext
	Fn             HandleFunc
	RequiredGroups []string
	RequiredPri    []string
	Owner          bool
}

func (h *BasicMngrHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	if h.RequiredGroups != nil || h.RequiredPri != nil || h.Owner {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
		user, err := h.Users.Get(token)
		if err != nil {
			if err == model.ErrNotLogged {
				JSONError(rw, err.Error(), http.StatusForbidden)
				return
			}

			JSONError(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// check for the current user
		if h.Owner {
			if EqualIdChecker(user.Id, mux.Vars(req)["user_id"]) {
				goto NORMAL
			}
			JSONError(rw, ErrForbidden.Error(), http.StatusForbidden)
			return
		}

		// check if any groups of the current user match one of the required groups
		if len(h.RequiredGroups) > 0 {
			for _, bg := range user.BriefGroups {
				for _, g2 := range h.RequiredGroups {
					if *bg.Name == g2 {
						goto NORMAL
					}
				}
			}
		}

		// check if any privileges of the current user match one of the required privileges
		if len(h.RequiredPri) > 0 {
			for _, pri := range user.Privilege {
				for _, p := range h.RequiredPri {
					if pri == p {
						goto NORMAL
					}
				}
			}
		}

		// check if any groups of the current user has the privileges match one of required privileges
		aid := make([]interface{}, 0, len(user.BriefGroups))
		for _, v := range user.BriefGroups {
			aid = append(aid, v.Id)
		}

		groups, err := h.AuthContext.Groups.FindSome(aid...)
		if err != nil {
			JSONError(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, v := range groups {
			for _, pri := range v.Privilege {
				for _, p := range h.RequiredPri {
					if pri == p {
						goto NORMAL
					}
				}
			}
		}

		JSONError(rw, err.Error(), http.StatusForbidden)
		return
	}

NORMAL:
	status, err := h.Fn(&h.AuthContext, rw, req)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		JSONError(rw, err.Error(), status)
	} else {
		rw.WriteHeader(status)
	}
}

func JSONError(rw http.ResponseWriter, message string, code int) {
	rw.WriteHeader(code)
	rw.Write([]byte(`{"error":"` + message + `"}`))
}
