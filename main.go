package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

func main() {
	// set up logging
	logger.Init()

	// load the config
	helpers.LoadConfig()

	// init hydra
	err := helpers.InitHydra()
	if err != nil {
		log.Fatal(err)
	}

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
	logger.Info.Println("X-IDP server is running on " + helpers.Config.Host + " press CTRL-C to quit")
	err = http.ListenAndServe(helpers.Config.Host, CSRF(router))
	if err != nil {
		log.Fatal(err)
	}
}
