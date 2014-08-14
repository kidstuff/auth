package auth

import (
	"code.google.com/p/go.net/context"
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/model"
	"github.com/kidstuff/conf"
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
	if HANDLER_REGISTER == nil {
		panic("kidstuff/auth: HANDLER_REGISTER need to be overide by a mngr")
	}

	if DEFAULT_NOTIFICATOR == nil {
		panic("kidstuff/auth: DEFAULT_NOTIFICATOR need to be overide by a mngr")
	}

	if ID_FROM_STRING == nil {
		panic("kidstuff/auth: ID_FROM_STRING need to be overide by a mngr")
	}

	if ID_TO_STRING == nil {
		panic("kidstuff/auth: ID_TO_STRING need to be overide by a mngr")
	}

	router.Handle("/signup", HANDLER_REGISTER(SignUp, false, nil, nil))
	router.Handle("/tokens",
		HANDLER_REGISTER(GetToken, false, nil, nil))

	router.Handle("/users/{user_id}/activate",
		HANDLER_REGISTER(Activate, false, nil, nil))

	router.Handle("/users/{user_id}/password",
		HANDLER_REGISTER(UpdatePassword, true, []string{"admin"}, []string{"manage_user"})).Methods("PUT")

	router.Handle("/users/{user_id}",
		HANDLER_REGISTER(GetProfile, false, nil, nil)).Methods("GET")

	router.Handle("/users/{user_id}",
		HANDLER_REGISTER(UpdateProfile, true, []string{"admin"}, []string{"manage_user"})).Methods("PATCH")

	router.Handle("/users",
		HANDLER_REGISTER(ListProfile, false, nil, nil))

}

type AuthContext struct {
	context.Context
	Users         model.UserManager
	Groups        model.GroupManager
	Settings      conf.Configurator
	Notifications Notificator
	Logs          Logger
}

type HandleFunc func(*AuthContext, http.ResponseWriter, *http.Request) (int, error)

var (
	HANDLER_REGISTER func(fn HandleFunc, owner bool, groups, pri []string) http.Handler
	ID_FROM_STRING   func(string) (interface{}, error)
	ID_TO_STRING     func(interface{}) (string, error)
)

type Condition struct {
	RequiredGroups []string
	RequiredPri    []string
	Owner          bool
}

func BasicMngrHandler(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request, cond *Condition, fn HandleFunc) {
	var cancel context.CancelFunc
	authCtx.Context, cancel = context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	authCtx.Notifications = DEFAULT_NOTIFICATOR
	authCtx.Logs, _ = NewSysLogger("kidstuff/auth")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	if cond.RequiredGroups != nil || cond.RequiredPri != nil || cond.Owner {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
		user, err := authCtx.Users.Get(token)
		if err != nil {
			if err == model.ErrNotLogged {
				JSONError(rw, err.Error(), http.StatusForbidden)
				return
			}

			JSONError(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// check for the current user
		if cond.Owner {
			if sid, _ := ID_TO_STRING(user.Id); sid == mux.Vars(req)["user_id"] {
				goto NORMAL
			}

			JSONError(rw, ErrForbidden.Error(), http.StatusForbidden)
			return
		}

		// check if any groups of the current user match one of the required groups
		if len(cond.RequiredGroups) > 0 {
			for _, bg := range user.BriefGroups {
				for _, g2 := range cond.RequiredGroups {
					if *bg.Name == g2 {
						goto NORMAL
					}
				}
			}
		}

		// check if any privileges of the current user match one of the required privileges
		if len(cond.RequiredPri) > 0 {
			for _, pri := range user.Privilege {
				for _, p := range cond.RequiredPri {
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

		groups, err := authCtx.Groups.FindSome(aid...)
		if err != nil {
			JSONError(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, v := range groups {
			for _, pri := range v.Privilege {
				for _, p := range cond.RequiredPri {
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
	status, err := fn(authCtx, rw, req)
	if err != nil {
		authCtx.Logs.Errorf("HTTP %d: %q", status, err)
		JSONError(rw, err.Error(), status)
	}
}

func JSONError(rw http.ResponseWriter, message string, code int) {
	rw.WriteHeader(code)
	rw.Write([]byte(`{"error":"` + message + `"}`))
}
