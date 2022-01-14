package auth

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
	"x-net.at/idp/authenticators"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

func Login(w http.ResponseWriter, request *http.Request) {
	// get the hydra admin
	hydraAdmin := helpers.GetAdmin()

	// get the requests context
	ctx := request.Context()

	// reject non post requests
	if request.Method != http.MethodPost {
		helpers.Error(w, http.StatusMethodNotAllowed, "Error: This endpoint only accepts POST requests!")
		return
	}

	// get the authenticator from mux
	vars := mux.Vars(request)
	realm := vars["realm"]

	// parse the form
	err := request.ParseForm()
	if err != nil {
		helpers.Error(w, 500, "Failed to parse login form: "+err.Error())
		return
	}

	// convert the remember me form data to bool
	var rememberMe bool = false

	profile, authOk := authenticators.Login(request.FormValue("email"), request.FormValue("password"), realm)
	if !authOk {
		logger.Log.Info("login failed")
		helpers.Render(w, request.Header.Get("Accept-Language"), "login.html", "base.html", helpers.TemplateCtx{Controller: LoginData{csrf.TemplateField(request), request.FormValue("login-challenge"), true, "username or password is wrong"}})
		return
	}

	// set remember in context
	profile.Remember = rememberMe
	logger.Log.Debug("profile", profile)

	//get the login challenge
	challenge := request.FormValue("login-challenge")
	// accept the login request with hydra admin
	loginGetParam := admin.NewGetLoginRequestParams()
	loginGetParam.SetLoginChallenge(challenge)

	_, err = hydraAdmin.GetLoginRequest(loginGetParam)
	if err != nil {
		helpers.Error(w, 500, "Failed to get login request: "+err.Error())
		return
	}

	subject := profile.Email

	loginAcceptParam := admin.NewAcceptLoginRequestParams()
	loginAcceptParam.WithContext(ctx)

	loginAcceptParam.SetLoginChallenge(challenge)
	loginAcceptParam.SetBody(&models.AcceptLoginRequest{
		Subject:  &subject,
		Remember: rememberMe,
		Context:  profile,
	})

	respLoginAccept, err := hydraAdmin.AcceptLoginRequest(loginAcceptParam)
	if err != nil {
		helpers.Error(w, 500, "Failed to accept login request: "+err.Error())
		return
	}

	// if the request is accepted redirect to the consent handler
	http.Redirect(w, request, *respLoginAccept.GetPayload().RedirectTo, http.StatusFound)
	return
}
