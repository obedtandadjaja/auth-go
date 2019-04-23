package main

import (
	"net/http"
	"github.com/obedtandadjaja/auth-go/controller"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Token",
		"POST",
		"/token",
		http.handlerFunc(controller.Token),
	},
}
