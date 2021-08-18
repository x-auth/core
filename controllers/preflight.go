package controllers

import (
	"net/http"
	"strings"
	"log"
	"x-net.at/idp/helpers"
)

/*
	This controller checks the post request of the form and decides which realm (and therefore authenticator) to use,
	by redirecting to the correct auth endpoint.
*/

func loginRedir(w http.ResponseWriter, request *http.Request, realm string) {
	http.Redirect(w, request, helpers.Config.BasePath+"/login/"+realm, http.StatusTemporaryRedirect)
}

func Preflight(w http.ResponseWriter, request *http.Request) {
	// reject non post requests
	if request.Method != http.MethodPost {
		helpers.Error(w, http.StatusMethodNotAllowed, "Error: This endpoint only accepts POST requests!")
		return
	}

	err := request.ParseForm()
	if err != nil {
		helpers.Error(w, 500, "Failed to parse login form: "+err.Error())
		return
	}

	ok, splitChar := helpers.SliceContains(request.FormValue("email"), helpers.Config.SplitCharacters)

	if !ok {
		// no valid split char, try the default realm
		for _, realm := range helpers.Config.Realms {
			if realm.Default {
				// redirect to the default authenticator login
				loginRedir(w, request, realm.Identifier)
				return
			}
		}
		helpers.Error(w, 400, "No default realm found")
		return
	} else {
		// split char is good
		identifierRealm := strings.Split(request.FormValue("email"), splitChar)[1]

		for _, realm := range helpers.Config.Realms {
			// check if the realm in the identifier matches a configured realm
			if realm.Identifier == identifierRealm {
				for _, authCfg := range helpers.Config.Authenticators {
					if authCfg.Name == realm.Authenticator {
						// redirect to the evaluated authenticator
						loginRedir(w, request, realm.Name)
						return
					}
				}
			}
		}
	}
	helpers.Error(w, 400, "No valid realm found")
}
