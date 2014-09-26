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
	OnlineThreshold  = time.Hour
	HANDLER_REGISTER func(fn HandleFunc, owner bool, groups, pri []string) http.Handler
	ID_FROM_STRING   func(string) (interface{}, error)
	ID_TO_STRING     func(interface{}) (string, error)
)

type ctxKey int

const (
	userTokenKey ctxKey = iota
	userIdKey
)

type HandleFunc func(*AuthContext, http.ResponseWriter, *http.Request) (int, error)

type AuthContext struct {
	context.Context
	Users         model.UserManager
	Groups        model.GroupManager
	Settings      conf.Configurator
	Notifications Notificator
	Logs          Logger
	currentUser   *model.User
}

func (ctx *AuthContext) saveToken(token string) {
	ctx.Context = context.WithValue(ctx.Context, userTokenKey, token)
}

func (ctx *AuthContext) saveId(id string) {
	ctx.Context = context.WithValue(ctx.Context, userIdKey, id)
}

func (ctx *AuthContext) ValidCurrentUser(owner bool, groups, pri []string) (*model.User, error) {
	if ctx.currentUser == nil {
		//try to query current user
		token, ok := ctx.Value(userTokenKey).(string)
		if !ok || len(token) == 0 {
			return nil, ErrForbidden
		}
		var err error
		ctx.currentUser, err = ctx.Users.Get(token)
		if err != nil {
			return nil, err
		}
		// calculate user privilege base on user's privilege and group's privilege
		mPri := make(map[string]bool)
		for _, p := range ctx.currentUser.Privilege {
			mPri[p] = true
		}

		aid := make([]interface{}, 0, len(ctx.currentUser.Groups))
		for _, v := range ctx.currentUser.Groups {
			aid = append(aid, v.Id)
		}

		groups, err := ctx.Groups.FindSome(aid...)
		if err == nil {
			for _, v := range groups {
				for _, p := range v.Privilege {
					mPri[p] = true
				}
			}
		} else {
			ctx.Logs.Errorf("cannot load user groups to get group privilege")
		}

		aPri := make([]string, 0, len(mPri))
		for p := range mPri {
			aPri = append(aPri, p)
		}

		ctx.currentUser.Privilege = aPri
	}

	err := validCurrentUser(ctx, ctx.currentUser, owner, groups, pri)
	return ctx.currentUser, err
}

func validCurrentUser(authCtx *AuthContext, user *model.User, owner bool, groups, privilege []string) error {
	// check for the current user
	if owner {
		sid, ok := authCtx.Context.Value(userIdKey).(string)
		uid, _ := ID_TO_STRING(user.Id)
		if !ok || len(sid) == 0 || sid != uid {
			return ErrForbidden
		}
	}

	// check if any groups of the current user match one of the required groups
	if len(groups) > 0 {
		foundGroup := false
	LOOP_GROUP:
		for _, bg := range user.Groups {
			for _, g2 := range groups {
				if *bg.Name == g2 {
					foundGroup = true
					break LOOP_GROUP
				}
			}
		}

		if !foundGroup {
			return ErrForbidden
		}
	}

	// check if any privileges of the current user match one of the required privileges
	if len(privilege) > 0 {
		foundPri := false
	LOOP_PRI:
		for _, pri := range user.Privilege {
			for _, p := range privilege {
				if pri == p {
					foundPri = true
					break LOOP_PRI
				}
			}
		}

		if !foundPri {
			return ErrForbidden
		}
	}

	return nil
}

type Condition struct {
	RequiredGroups []string
	RequiredPri    []string
	Owner          bool
}

func BasicMngrHandler(authCtx *AuthContext, rw http.ResponseWriter, req *http.Request, cond *Condition, fn HandleFunc) {
	var cancel context.CancelFunc
	authCtx.Context, cancel = context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	authCtx.saveToken(token)
	authCtx.saveId(mux.Vars(req)["user_id"])

	authCtx.Notifications = DEFAULT_NOTIFICATOR
	authCtx.Logs, _ = NewSysLogger("kidstuff/auth")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	if cond.RequiredGroups != nil || cond.RequiredPri != nil || cond.Owner {
		_, err := authCtx.ValidCurrentUser(cond.Owner, cond.RequiredGroups, cond.RequiredPri)
		if err != nil {
			JSONError(rw, err.Error(), http.StatusForbidden)
			return
		}
	}

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
