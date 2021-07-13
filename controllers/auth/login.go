package auth

import (
	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"html/template"
	"net/http"
	"x-net.at/idp/authenticators"
	"x-net.at/idp/helpers"
)

type LoginData struct {
	CSRFField template.HTML
	Challenge string
	Error     bool
	Message   string
}

func Login(w http.ResponseWriter, request *http.Request) {
	hydraAdmin := helpers.GetAdmin()
	if request.Method == http.MethodPost {
		// POST Handler
		// get the context
		ctx := request.Context()
		defer ctx.Done()

		// parse the form
		err := request.ParseForm()
		if err != nil {
			helpers.Error(w, 500, "Failed to parse login form: "+err.Error())
			return
		}

		// convert the remember me form data to bool
		var rememberMe bool = false
		if request.FormValue("remember") == "true" {
			rememberMe = true
		}

		// validate username and password by trying all authenticators
		cookie, err := request.Cookie("x-idp-authenticator")
		if err != nil {
			helpers.Error(w, 500, "Cookie Error: "+err.Error())
		}

		profile, authOK := authenticators.Login(request.FormValue("email"), request.FormValue("password"), cookie, &w)

		// authentication failed, re-render the login form with error
		if !authOK {
			helpers.Render(w, "login.html", "base.html", LoginData{csrf.TemplateField(request), request.FormValue("login-challenge"), true, "username or password is wrong"})
			return
		}

		// use Hydra admin to accept the login request
		loginGetParam := admin.NewGetLoginRequestParams()
		loginGetParam.SetLoginChallenge(request.FormValue("login-challenge"))

		_, err = hydraAdmin.GetLoginRequest(loginGetParam)
		if err != nil {
			helpers.Error(w, 500, "Failed to get login request: "+err.Error())
			return
		}

		subject := profile.Email
		loginAcceptParam := admin.NewAcceptLoginRequestParams()
		loginAcceptParam.WithContext(ctx)
		loginAcceptParam.SetLoginChallenge(request.FormValue("login-challenge"))
		loginAcceptParam.SetBody(&models.AcceptLoginRequest{
			Subject:  &subject,
			Remember: rememberMe,
		})

		respLoginAccept, err := hydraAdmin.AcceptLoginRequest(loginAcceptParam)
		if err != nil {
			helpers.Error(w, 500, "Failed to accept login request: "+err.Error())
			return
		}

		// if the request is accepted redirect to the consent handler
		http.Redirect(w, request, *respLoginAccept.GetPayload().RedirectTo, http.StatusFound)
		return
	} else {
		//GET Handler
		// get context
		ctx := request.Context()
		defer ctx.Done()

		// get the login challenge
		challenge_slice, ok := request.URL.Query()["login_challenge"]
		if !ok || len(challenge_slice) < 1 {
			helpers.Error(w, 400, "Expected a login challenge but received none")
			return
		}

		// use Hydra Admin to get the login challenge info
		loginGetParam := admin.NewGetLoginRequestParams()
		loginGetParam.WithContext(ctx)
		loginGetParam.SetLoginChallenge(challenge_slice[0])

		respLoginGet, err := hydraAdmin.GetLoginRequest(loginGetParam)
		if err != nil {
			helpers.Error(w, 500, "Failed to get login request info: "+err.Error())
			return
		}

		// convert bool pointer to bool
		skip := false
		if respLoginGet.GetPayload().Skip != nil {
			skip = *respLoginGet.GetPayload().Skip
		}

		// If hydra was already able to authenticate the user, skip will be true and we do not need to re-authenticate
		// the user.
		if skip {
			// Use Hydra admin to accept the login request
			loginAcceptParam := admin.NewAcceptLoginRequestParams()
			loginAcceptParam.WithContext(ctx)
			loginAcceptParam.SetLoginChallenge(challenge_slice[0])
			loginAcceptParam.SetBody(&models.AcceptLoginRequest{
				Subject: respLoginGet.GetPayload().Subject,
			})

			respLoginAccept, err := hydraAdmin.AcceptLoginRequest(loginAcceptParam)
			if err != nil {
				helpers.Error(w, 400, "Cannot accept login request: "+err.Error())
				return
			}

			http.Redirect(w, request, *respLoginAccept.GetPayload().RedirectTo, http.StatusFound)
			return
		}
		helpers.Render(w, "login.html", "base.html", LoginData{csrf.TemplateField(request), challenge_slice[0], false, ""})
	}
}
