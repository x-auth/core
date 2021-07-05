package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"x-net.at/idp/helpers"
)

func main() {
	// load the config
	helpers.LoadConfig()

	// init hydra
	err := helpers.InitHydra()
	if err != nil {
		log.Fatal(err)
	}

	// init kratos
	helpers.InitKratosClient()

	// handle the urls
	router := mux.NewRouter()
	for _, route := range routes {
		router.HandleFunc(route.path, route.handler)
	}

	// handle staticfiles
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	//load secret key
	helpers.LoadSecretKey()

	// set up csrf protection
	CSRF := csrf.Protect(helpers.SecretKey)

	//init secure cookie
	helpers.InitSecureCookie()

	// start the server
	log.Println("Auth server is running on http://localhost:8000 press CTRL-C to quit")
	err = http.ListenAndServe(":8000", CSRF(router))
	if err != nil {
		log.Fatal(err)
	}
}
