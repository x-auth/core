package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"nictec.net/auth/controllers"
	"nictec.net/auth/helpers"
)

func main(){
	// load the config
	helpers.LoadConfig()

	// init hydra
	controllers.InitHydra()

	// handle the urls
	router := mux.NewRouter()
	for _, route := range routes{
		router.HandleFunc(route.path, route.handler)
	}

	// handle staticfiles
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	//set up csrf protection
	key, err := ioutil.ReadFile("system/secret.key")
	if err != nil{
		log.Fatal(err)
	}
	CSRF := csrf.Protect(key)

	// start the server
	log.Println("Auth server is running on http://localhost:8000 press CTRL-C to quit")
	err = http.ListenAndServe(":8000", CSRF(router))
	if err != nil{
		log.Fatal(err)
	}
}
