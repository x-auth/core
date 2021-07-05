package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"x-net.at/idp/helpers"
)

type authCheckRequest struct {
	Identifier string `json:"identifier"`
}

type authCheckResponse struct {
	NeedsRedirect bool `json:"needs_redirect"`
}

func GetAuthenticator(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		helpers.Error(w, http.StatusMethodNotAllowed, "Error: This endpoint only accepts POST requests!")
		return
	}

	var req authCheckRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		helpers.JsonError(w, 400, "Decode error: "+err.Error())
		return
	}

	identifier := req.Identifier

	// check if a valid split char is in the identifier
	ok, splitChar := helpers.SliceContains(identifier, helpers.Config.SplitCharacters)
	if !ok {
		for _, realm := range helpers.Config.Realms {
			if realm.Default {
				for _, authenticator := range helpers.Config.Authenticators {
					if authenticator.Type == "kratos" {
						// tell client to post to kratos
						helpers.JsonResponse(w, authCheckResponse{NeedsRedirect: true})
						return
					} else {
						// set the authenticator cookie, so that we dont have to evaluate the config twice
						value := map[string]string{
							"authenticator": authenticator.Type,
							"realm":         realm.Identifier,
						}
						encoded, err := helpers.SecureCookie.Encode("x-idp-authenticator", value)
						if err != nil {
							helpers.JsonError(w, 500, "Cookie Error: "+err.Error())
							return
						}
						http.SetCookie(w, &http.Cookie{Name: "x-idp-authenticator", Value: encoded, Secure: true, HttpOnly: true})

						// respond that no "redirect" is needed
						helpers.JsonResponse(w, authCheckResponse{NeedsRedirect: false})
						return
					}
				}
			}
		}
	}
	// ok
	realmName := strings.Split(identifier, splitChar)[1]
	for _, realm := range helpers.Config.Realms {
		// check if the realm in the identifier matches a configured realm
		if realm.Identifier == realmName {
			for _, authCfg := range helpers.Config.Authenticators {
				if authCfg.Name == realm.Authenticator {
					if authCfg.Type == "kratos" {
						helpers.JsonResponse(w, authCheckResponse{NeedsRedirect: true})
						return
					} else {
						// set the authenticator cookie, so that we dont have to evaluate the config twice
						value := map[string]string{
							"authenticator": authCfg.Type,
							"realm":         realm.Identifier,
						}
						encoded, err := helpers.SecureCookie.Encode("x-idp-authenticator", value)
						if err != nil {
							helpers.JsonError(w, 500, "Cookie Error: "+err.Error())
							return
						}
						http.SetCookie(w, &http.Cookie{Name: "x-idp-authenticator", Value: encoded, Secure: true, HttpOnly: true})

						// tell the client to not redirect
						helpers.JsonResponse(w, authCheckResponse{NeedsRedirect: false})
						return
					}
				}
			}
		}
	}
	helpers.JsonError(w, 400, "Unknown realm")
}

func Urls(w http.ResponseWriter, request *http.Request) {
	urls := struct {
		HydraURL  string `json:"hydra_url"`
		KratosURL string `json:"kratos_url"`
	}{helpers.Config.HydraURL, helpers.Config.KratosURL}
	helpers.JsonResponse(w, urls)
}
