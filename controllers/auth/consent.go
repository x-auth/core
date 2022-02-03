package auth

import (
	"context"
	"encoding/json"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

func Consent(w http.ResponseWriter, request *http.Request) {
	// reject POST request
	if request.Method == http.MethodPost {
		helpers.Error(w, 400, "Post requests are not allowed for this endpoint")
		return
	}

	hydraAdmin := helpers.GetAdmin()
	ctx := request.Context()
	challenge_slice, ok := request.URL.Query()["consent_challenge"]
	if !ok || len(challenge_slice) < 1 {
		logger.Log.Warning("Expected a login challenge but received none")
		helpers.Error(w, 400, "Expected a login challenge but received none")
		return
	}

	// use the hydra admin to get the consent challenge info
	consentGetParams := admin.NewGetConsentRequestParams()
	consentGetParams.WithContext(ctx)
	consentGetParams.SetConsentChallenge(challenge_slice[0])

	// get the consent data from hydra
	consentGetResp, err := hydraAdmin.GetConsentRequest(consentGetParams)
	if err != nil {
		logger.Log.Error("Failed to get consent request, " + err.Error())
		helpers.Error(w, 500, "Failed to get consent request: "+err.Error())
		return
	}

	// get the requested grant scope
	grantScope := consentGetResp.GetPayload().RequestedScope

	// handle the session
	session := &models.ConsentRequestSession{}
	profile := consentGetResp.GetPayload().Context
	claims := helpers.GetClaims(grantScope)

	// compile time cast check, because we dont want runtime panics
	var parsedProfile map[string]interface{}
	switch t := profile.(type) {
	default:
		// TODO: Handle fatal error
	case map[string]interface{}:
		parsedProfile = t
	}

	if len(parsedProfile) == 0 {
		redisCtx := context.Background()
		logger.Log.Debug(consentGetResp.GetPayload().Subject)
		val, err := helpers.RDB.Get(redisCtx, consentGetResp.GetPayload().Subject).Result()
		err = json.Unmarshal([]byte(val), &parsedProfile)
		if err != nil {
			helpers.Error(w, 500, err.Error())
			return
		}
	}

	rememberMe := parsedProfile["Remember"].(bool)

	IDToken := make(map[string]interface{})
	for _, claim := range claims {
		IDToken[helpers.ToSnakeCase(claim)] = parsedProfile[claim].(string)
	}
	
	IDToken["email_verified"] = true
	session.IDToken = IDToken

	if helpers.Contains(consentGetResp.GetPayload().RequestedScope, "email") {
		grantScope = append(grantScope, "openid")
	}

	consentAcceptBody := &models.AcceptConsentRequest{
		GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
		GrantScope:               grantScope,
		Remember:                 rememberMe,
		Session:                  session,
	}

	consentAcceptParams := admin.NewAcceptConsentRequestParams()
	consentAcceptParams.WithContext(ctx)
	consentAcceptParams.SetConsentChallenge(challenge_slice[0])
	consentAcceptParams.WithBody(consentAcceptBody)

	consentAcceptResp, err := hydraAdmin.AcceptConsentRequest(consentAcceptParams)
	if err != nil {
		logger.Log.Error(err)
		helpers.Error(w, 500, "Failed to accept consent Request: "+err.Error())
		return
	}

	http.Redirect(w, request, *consentAcceptResp.GetPayload().RedirectTo, http.StatusFound)
}
