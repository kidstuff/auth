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

func Serve(router *mux.Router) {
	if HandlerRegister == nil {
		panic("kidstuff/auth: HandlerRegister need to be initialed by a mngr")
	}

	router.Handle("/tokens", HandlerRegister(GetToken, false, nil, nil))

}

type AuthContext struct {
	Users       model.UserManager
	Groups      model.GroupManager
	Settings    model.Configurator
	currentUser *model.User
}

func (ctx *AuthContext) CurrentUser(token string) (*model.User, error) {
	var err error
	if ctx.currentUser == nil {
		ctx.currentUser, err = ctx.Users.Get(token)
	}

	return ctx.currentUser, err
}

type HandleFunc func(*AuthContext, http.ResponseWriter, *http.Request) (int, error)

var HandlerRegister func(fn HandleFunc, owner bool, groups, pri []string) http.Handler

type BasicMngrHandler struct {
	AuthContext
	Fn             HandleFunc
	RequiredGroups []string
	RequiredPri    []string
	Owner          bool
}

func (h *BasicMngrHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	if h.RequiredGroups != nil || h.RequiredPri != nil {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
		user, err := h.CurrentUser(token)
		if err != nil {
			if err == model.ErrNotLogged {
				JSONError(rw, err.Error(), http.StatusForbidden)
				return
			}

			JSONError(rw, err.Error(), http.StatusInternalServerError)
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
	}
}

func JSONError(rw http.ResponseWriter, message string, code int) {
	rw.WriteHeader(code)
	rw.Write([]byte(`{"error":"` + message + `"}`))
}
