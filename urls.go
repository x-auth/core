package main

import (
	"net/http"
	"nictec.net/auth/controllers"
)

type path struct {
	path string
	handler func(w http.ResponseWriter, r *http.Request)
}

var routes []path = []path{
	{"/", controllers.Login},
	//{"/logout", controllers.Logout},
	{"/consent", controllers.Consent},
}