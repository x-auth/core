package controllers

import (
	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"html/template"
	"net/http"
	"x-net.at/idp/authenticators"
	"x-net.at/idp/helpers"
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
			helpers.Error(w, 500, "Failed to get consent request: "+err.Error())
			return
		}

		// grant the consent request
		var grantScope = request.Form["grant_scope"]

		// parse the remember form value to boolean
		remember := (request.FormValue("remember") == "true")

		// handle the session
		session := &models.ConsentRequestSession{}
		if helpers.Contains(grantScope, "profile") {
			cookie, err := request.Cookie("x-idp-authenticator")
			if err != nil {
				helpers.Error(w, 500, "cookie error: "+err.Error())
			}

			value := make(map[string]string)
			err = helpers.SecureCookie.Decode("x-idp-authenticator", cookie.Value, &value)
			if err != nil {
				helpers.Error(w, 500, "cookie error: "+err.Error())
			}

			profile := authenticators.GetProfile(value["authenticator"], consentGetResp.GetPayload().Subject)
			session.IDToken = profile
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
			helpers.Error(w, 500, "Failed to accept consent Request: "+err.Error())
			return
		}

		http.Redirect(w, request, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)
	} else {
		// GET Handler
		// get the context
		ctx := request.Context()
		defer ctx.Done()

		// get the login challenge
		challenge_slice, ok := request.URL.Query()["consent_challenge"]
		if !ok || len(challenge_slice) < 1 {
			helpers.Error(w, 400, "Expected a login challenge but received none")
			return
		}

		// use the Hydra admin to get consent challenge info
		consentGetParams := admin.NewGetConsentRequestParams()
		consentGetParams.WithContext(ctx)
		consentGetParams.SetConsentChallenge(challenge_slice[0])

		consentGetResp, err := hydraAdmin.GetConsentRequest(consentGetParams)
		if err != nil {
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
				helpers.Error(w, 500, "Failed to accept consent Request: "+err.Error())
				return
			}

			http.Redirect(w, request, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)
			return
		}
		// show the consent page
		type ConsentData struct {
			CSRFField      template.HTML
			Challenge      string
			RequestedScope models.StringSlicePipeDelimiter
			User           string
			Client         string
		}
		helpers.Render(w, "consent.html", "base.html", ConsentData{csrf.TemplateField(request), challenge_slice[0], consentGetResp.GetPayload().RequestedScope, consentGetResp.GetPayload().Subject, consentGetResp.GetPayload().Client.ClientName})
	}
}
