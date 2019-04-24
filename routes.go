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
		"CreateCredential",
		"POST",
		"/credentials",
		credentials.Create,
	},
}
