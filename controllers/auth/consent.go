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

package auth

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"html/template"
	"net/http"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

func Consent(w http.ResponseWriter, request *http.Request) {
	hydraAdmin := helpers.GetAdmin()
	if request.Method == http.MethodPost {
		// POST Handler
		// get the context
		ctx := request.Context()
		defer ctx.Done()

		// parse the form
		request.ParseForm()

		// use the Hydra admin to get consent challenge info
		consentGetParams := admin.NewGetConsentRequestParams()
		consentGetParams.WithContext(ctx)
		consentGetParams.SetConsentChallenge(request.FormValue("consent-challenge"))

		consentGetResp, err := hydraAdmin.GetConsentRequest(consentGetParams)
		if err != nil {
			logger.Error.Println("Failed to get consent request, " + err.Error())
			helpers.Error(w, 500, "Failed to get consent request: "+err.Error())
			return
		}

		// grant the consent request
		var grantScope = request.Form["grant_scope"]

		// parse the remember form value to boolean
		remember := request.FormValue("remember") == "true"

		// handle the session
		session := &models.ConsentRequestSession{}
		profile := consentGetResp.GetPayload().Context

		// handle the case that login was skipped
		if len(profile.(map[string]interface{})) == 0 {
			fmt.Println("old claims:")
			fmt.Println(consentGetResp.GetPayload().OidcContext)
		}

		claims := helpers.GetClaims(grantScope)

		// compile time cast check because go does not provide a good way for that
		var parsedProfile map[string]interface{}
		switch t := profile.(type) {
		default:
			// TODO: Handle fatal error
		case map[string]interface{}:
			parsedProfile = t
		}

		var IDToken = make(map[string]interface{})
		for _, claim := range claims {
			IDToken[helpers.ToSnakeCase(claim)] = parsedProfile[claim].(string)
		}
		session.IDToken = IDToken

		if helpers.Contains(consentGetResp.GetPayload().RequestedScope, "email") {
			IDToken["email_verified"] = true
		}

		// always append the openid grant scope
		if helpers.Contains(consentGetResp.GetPayload().RequestedScope, "openid") {
			grantScope = append(grantScope, "openid")
		}

		consentAcceptBody := &models.AcceptConsentRequest{
			GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
			GrantScope:               grantScope,
			Remember:                 remember,
			Session:                  session,
		}

		consentAcceptParams := admin.NewAcceptConsentRequestParams()
		consentAcceptParams.WithContext(ctx)
		consentAcceptParams.SetConsentChallenge(request.FormValue("consent-challenge"))
		consentAcceptParams.WithBody(consentAcceptBody)

		consentAcceptResp, err := hydraAdmin.AcceptConsentRequest(consentAcceptParams)
		if err != nil {
			logger.Error.Println(err)
			helpers.Error(w, 500, "Failed to accept consent Request: "+err.Error())
			return
		}

		http.Redirect(w, request, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)
	} else {
		// GET Handler
		// get the language
		lang := request.Header.Get("Accept-Language")
		// get the context
		ctx := request.Context()
		defer ctx.Done()

		// get the login challenge
		challenge_slice, ok := request.URL.Query()["consent_challenge"]
		if !ok || len(challenge_slice) < 1 {
			logger.Warning.Println("Expected a login challenge but received none")
			helpers.Error(w, 400, "Expected a login challenge but received none")
			return
		}

		// use the Hydra admin to get consent challenge info
		consentGetParams := admin.NewGetConsentRequestParams()
		consentGetParams.WithContext(ctx)
		consentGetParams.SetConsentChallenge(challenge_slice[0])

		consentGetResp, err := hydraAdmin.GetConsentRequest(consentGetParams)
		if err != nil {
			logger.Error.Println("Failed to get consent request: " + err.Error())
			helpers.Error(w, 500, "Failed to get consent request: "+err.Error())
			return
		}

		// If a user has granted this application the requested scope, hydra will tell us to not show the UI.
		if consentGetResp.GetPayload().Skip {
			// grant the consent request
			consentAcceptBody := &models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
				GrantScope:               consentGetResp.GetPayload().RequestedScope,
			}
			consentAcceptParams := admin.NewAcceptConsentRequestParams()
			consentAcceptParams.WithContext(ctx)
			consentAcceptParams.SetConsentChallenge(challenge_slice[0])
			consentAcceptParams.WithBody(consentAcceptBody)

			consentAcceptResp, err := hydraAdmin.AcceptConsentRequest(consentAcceptParams)
			if err != nil {
				logger.Error.Println("Failed to get consent request: " + err.Error())
				helpers.Error(w, 500, "Failed to accept consent Request: "+err.Error())
				return
			}

			http.Redirect(w, request, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)
			return
		}
		// show the consent page
		type ConsentData struct {
			CSRFField template.HTML
			Challenge string
			//TODO: remove openid from list of scopes to give consent and accept it by default
			//Note: Should we deny the flow if openid is not requested?
			RequestedScope map[string][]string
			User           string
			Client         string
		}

		var relevantScopes = make(map[string][]string)
		for _, scope := range consentGetResp.GetPayload().RequestedScope {
			if scope != "openid" {
				if lang == "de" {
					relevantScopes[scope] = helpers.ScopeClaimMapperGerman[scope]
				} else {
					relevantScopes[scope] = helpers.ScopeClaimMapperEnglish[scope]
				}
			}
		}

		helpers.Render(w, lang, "consent.html", "base.html", helpers.TemplateCtx{Controller: ConsentData{csrf.TemplateField(request), challenge_slice[0], relevantScopes, consentGetResp.GetPayload().Subject, consentGetResp.GetPayload().Client.ClientName}})
	}
}
