package auth

import (
	"github.com/gorilla/mux"
	"github.com/kidstuff/auth/model"
	"net/http"
)

func Serve(router *mux.Router) {

}

type appContext struct {
	Users  model.UserManager
	Groups model.GroupManager
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}
