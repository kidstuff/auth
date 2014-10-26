package auth

import (
	"github.com/gorilla/mux"
)

func Serve(router *mux.Router) {
	if HANDLER_REGISTER == nil {
		panic("kidstuff/auth: HANDLER_REGISTER need to be overide by a mngr")
	}

	if DEFAULT_NOTIFICATOR == nil {
		panic("kidstuff/auth: DEFAULT_NOTIFICATOR need to be overide by a mngr")
	}

	router.Handle("/signup", HANDLER_REGISTER(SignUp, false, nil, nil))
	router.Handle("/tokens",
		HANDLER_REGISTER(GetToken, false, nil, nil))

	router.Handle("/users/{user_id}/activate",
		HANDLER_REGISTER(Activate, false, nil, nil))

	router.Handle("/users/{user_id}/password",
		HANDLER_REGISTER(UpdatePassword, true, []string{"admin"}, []string{"manage_user"})).Methods("PUT")

	router.Handle("/users/{user_id}",
		HANDLER_REGISTER(GetUser, true, []string{"admin"}, []string{"manage_user"})).Methods("GET")

	router.Handle("/users/{user_id}/profile",
		HANDLER_REGISTER(UpdateUserProfile, true, []string{"admin"}, []string{"manage_user"})).Methods("PATCH")

	router.Handle("/users/{user_id}/approve",
		HANDLER_REGISTER(UpdateApprovedStatus, false, []string{"admin"}, []string{"manage_user"})).Methods("PUT")

	router.Handle("/users",
		HANDLER_REGISTER(ListUser, false, []string{"admin"}, []string{"manage_user"}))
}
