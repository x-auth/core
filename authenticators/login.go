package authenticators

import (
	"net/http"
	"strings"
	"x-net.at/idp/helpers"
)

func getConfig(authenticator string, realmId string) map[string]string {
	for _, realm := range helpers.Config.Realms {
		if realm.Identifier == realmId {
			for _, auth := range helpers.Config.Authenticators {
				if auth.Type == authenticator {
					return helpers.ReduceConfig(auth.Config, realm.Config)
				}
			}
		}
	}
	return nil
}

func Login(identifier string, password string, cookie *http.Cookie) (Profile, bool) {
	// check if a valid split char is in the identifier
	ok, splitChar := helpers.SliceContains(identifier, helpers.Config.SplitCharacters)
	if !ok {
		return Profile{}, false
	}

	// split the identifier in username and realm
	idSlice := strings.Split(identifier, splitChar)
	username := idSlice[0]
	realmName := idSlice[1]

	// get the authenticator and realm from the preflight via secure cookie
	value := make(map[string]string)
	err := helpers.SecureCookie.Decode("x-idp-authenticator", cookie.Value, &value)
	if err != nil {
		return Profile{}, false
	}

	preflightAuthenticator := value["authenticator"]
	preflightRealm := value["realm"]

	// quit if the input is wrong
	if realmName != preflightRealm {
		return Profile{}, false
	}

	// authenticate using the right authenticator
	if preflightAuthenticator == "mock" {
		return mock(username, password, getConfig(preflightAuthenticator, preflightRealm))
	}
	return Profile{}, false
}
