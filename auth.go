package auth

import (
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/model"
	"net/http"
	"time"
)

var (
	OnlineThreshold = time.Hour
)

func Serve(router *mux.Router) {
	if HandlerRegister == nil {
		panic("kidstuff/auth: HandlerRegister need to be initialed by a mngr")
	}

	router.Handle("/tokens", HandlerRegister(GetToken, nil, nil))

}

type AuthContext struct {
	Users       model.UserManager
	Groups      model.GroupManager
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

var HandlerRegister func(fn HandleFunc, groups []string, pri []string) http.Handler
