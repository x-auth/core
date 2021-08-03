package auth

import (
	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"html/template"
	"net/http"
	"x-net.at/idp/helpers"
)

type LoginData struct {
	CSRFField template.HTML
	Challenge string
	Error     bool
	Message   string
}

func LoginForm(w http.ResponseWriter, request *http.Request) {
	hydraAdmin := helpers.GetAdmin()

	// reject non get requests
	if request.Method != http.MethodGet{
		helpers.Error(w, http.StatusMethodNotAllowed, "Error: This endpoint only accepts POST requests!")
		return
	}

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
