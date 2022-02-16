/*
 * Copyright (c) 2021 X-Net Services GmbH
 * Info: https://x-net.at
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import "C"
import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"x-net.at/idp/authenticators"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

func main() {
	// set up logging
	logger.Init()
	defer logger.Destroy()

	// load the config
	conf := helpers.LoadConfig()
	for _, realm := range conf.Realms {
		if !authenticators.ValidateRealm(realm.Authenticator) {
			logger.Log.Fatal("realm \"", realm.Name, "\" contains non existent authenticator \"", realm.Authenticator, "\"")
		}
	}

	// load the plugin config
	authenticators.Init()

	// set up redis
	helpers.InitRedis()

	// init hydra
	// TODO: Tell hydra that we only support scopes defined in mappers.go
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
	// TODO: Add https option for secure false
	CSRF := csrf.Protect(helpers.SecretKey, csrf.Secure(false))

	//init secure cookie
	helpers.InitSecureCookie()

	// start the server
	logger.Log.Info("X-IDP server is running on " + helpers.Config.Host)
	logger.Log.Debug("Press CTRL+C to quit")
	err = http.ListenAndServe(helpers.Config.Host, CSRF(router))
	if err != nil {
		log.Fatal(err)
	}
}
