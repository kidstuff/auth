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

	router.Handle("/signup", HANDLER_REGISTER(SignUp, false, nil))
	router.Handle("/tokens",
		HANDLER_REGISTER(GetToken, false, nil))

	router.Handle("/users/{user_id}/activate",
		HANDLER_REGISTER(Activate, false, nil))

	router.Handle("/users/{user_id}/password",
		HANDLER_REGISTER(ChangePassword, true, nil)).Methods("PUT")

	router.Handle("/users/{user_id}/password/override",
		HANDLER_REGISTER(OverridePassword, false, []string{"manage_user"})).Methods("PUT")

	router.Handle("/users/{user_id}",
		HANDLER_REGISTER(GetUser, true, []string{"manage_user"})).Methods("GET")

	router.Handle("/users/{user_id}",
		HANDLER_REGISTER(DeleteUser, true, []string{"manage_user"})).Methods("DELETE")

	router.Handle("/users/{user_id}/profile",
		HANDLER_REGISTER(UpdateUserProfile, true, []string{"manage_user"})).Methods("PATCH")

	router.Handle("/users/{user_id}/approve",
		HANDLER_REGISTER(UpdateApprovedStatus, false, []string{"manage_user"})).Methods("PUT")

	router.Handle("/users/{user_id}/groups/{group_id}",
		HANDLER_REGISTER(RemoveGroupFromUser, false, []string{"manage_user"})).Methods("DELETE")

	router.Handle("/users/{user_id}/groups",
		HANDLER_REGISTER(AddGroupToUser, false, []string{"manage_user"})).Methods("PUT")

	router.Handle("/users",
		HANDLER_REGISTER(ListUser, false, []string{"manage_user"}))

	router.Handle("/groups",
		HANDLER_REGISTER(CreateGroup, false, []string{"manage_user"})).Methods("POST")

	router.Handle("/groups",
		HANDLER_REGISTER(ListGroup, false, []string{"manage_user"}))
}
