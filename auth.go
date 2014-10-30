package auth

import (
	"code.google.com/p/go.net/context"
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/authmodel"
	"github.com/kidstuff/conf"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var (
	OnlineThreshold = time.Hour
	// HANDLER_REGISTER should be "overided" by the "manager". Implement of this function
	// must use the "or" logic for the conditions.
	HANDLER_REGISTER func(fn HandleFunc, owner bool, groups, pri []string) http.Handler
)

type ctxKey int

const (
	userTokenKey ctxKey = iota
	userIdKey
)

type HandleFunc func(*AuthContext, http.ResponseWriter, *http.Request) (int, error)

type AuthContext struct {
	context.Context
	Auth          authmodel.Manager
	Settings      conf.Configurator
	Notifications Notificator
	Logs          Logger
	currentUser   *authmodel.User
}

func (ctx *AuthContext) saveToken(token string) {
	ctx.Context = context.WithValue(ctx.Context, userTokenKey, token)
}

func (ctx *AuthContext) saveId(id string) {
	ctx.Context = context.WithValue(ctx.Context, userIdKey, id)
}

// ValidCurrentUser validate user privilege and cacuate user total privilege base on groups
func (ctx *AuthContext) ValidCurrentUser(owner bool, groups, pri []string) (*authmodel.User, error) {
	if ctx.currentUser == nil {
		//try to query current user
		token, ok := ctx.Value(userTokenKey).(string)
		if !ok || len(token) == 0 {
			return nil, ErrForbidden
		}
		var err error
		ctx.currentUser, err = ctx.Auth.GetUser(token)
		if err != nil {
			return nil, err
		}
		// calculate user privilege base on user's privilege and group's privilege
		mPri := make(map[string]bool)
		for _, p := range ctx.currentUser.Privileges {
			mPri[p] = true
		}

		aid := make([]string, 0, len(ctx.currentUser.Groups))
		for _, v := range ctx.currentUser.Groups {
			aid = append(aid, *v.Id)
		}

		groups, err := ctx.Auth.FindSomeGroup(aid...)
		if err == nil {
			for _, v := range groups {
				for _, p := range v.Privileges {
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

		ctx.currentUser.Privileges = aPri
	}

	err := validCurrentUser(ctx, ctx.currentUser, owner, groups, pri)
	return ctx.currentUser, err
}

func validCurrentUser(authCtx *AuthContext, user *authmodel.User, owner bool, groups, privilege []string) error {
	// check for the current user
	if owner {
		sid, ok := authCtx.Context.Value(userIdKey).(string)
		if !ok || len(sid) == 0 || sid != *user.Id {
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
		for _, pri := range user.Privileges {
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

// BasicMngrHandler can be use in "manager" ServeHTTP after initital required interface like
// authmodel.UserManager, authmodel.GroupManager, conf.Configurator...etc
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

// JSONError is a helper function to write json error message to http.ResponseWriter
func JSONError(rw http.ResponseWriter, message string, code int) {
	rw.WriteHeader(code)
	rw.Write([]byte(`{"error":"` + template.JSEscapeString(message) + `"}`))
}
