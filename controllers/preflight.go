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

package controllers

import (
	"net/http"
	"strings"
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
		helpers.Error(w, 400, "Identifier contains no split character")
		return
	}

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
	for _, realm := range helpers.Config.Realms {
		if realm.Default {
			loginRedir(w, request, realm.Name)
			return
		}
	}

	helpers.Error(w, 400, "no valid realm found")
}
