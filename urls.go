package main

import (
	"net/http"
	"x-net.at/idp/controllers"
	"x-net.at/idp/controllers/auth"
)

type path struct {
	path    string
	handler func(w http.ResponseWriter, r *http.Request)
}

var routes []path = []path{
	{"/", auth.LoginForm},
	{"/login/{realm}", auth.Login},
	{"/consent", auth.Consent},
	{"/preflight", controllers.Preflight},
}
