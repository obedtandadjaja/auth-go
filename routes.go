package main

import (
	"net/http"

	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/controller/credentials"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error
}

type Routes []Route

var routes = Routes{
	Route{
		"Token",
		"POST",
		"/token",
		controller.Token,
	},
	Route{
		"Token",
		"POST",
		"/refresh-token",
		controller.Login,
	},
	Route{
		"Verify",
		"POST",
		"/verify",
		controller.Verify,
	},
	Route{
		"CreateCredential",
		"POST",
		"/credentials",
		credentials.Create,
	},
	Route{
		"DeleteCredential",
		"DELETE",
		"/credentials",
		credentials.Delete,
	},
	Route{
		"ResetPassword",
		"POST",
		"/credentials/reset_password",
		credentials.ResetPassword,
	},
	Route{
		"InitiateResetPassword",
		"POST",
		"/credentials/initiate_password_reset",
		credentials.InitiatePasswordReset,
	},
}
