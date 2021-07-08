package authenticators

import (
	"log"
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

func Login(identifier string, password string, cookie *http.Cookie, w *http.ResponseWriter) (Profile, bool) {
	// check if a valid split char is in the identifier
	ok, splitChar := helpers.SliceContains(identifier, helpers.Config.SplitCharacters)
	if !ok {
		return Profile{}, false
	}

	// split the identifier in username and realm
	idSlice := strings.Split(identifier, splitChar)
	username := identifier
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
	} else if preflightAuthenticator == "ldap" {
		profile, ok := ldap(username, password, getConfig(preflightAuthenticator, preflightRealm))
		if !ok {
			return Profile{}, false
		}

		// set the profile cookie
		encoded, err := helpers.SecureCookie.Encode("x-idp-profile", profile)
		if err != nil {
			log.Println(err)
			return Profile{}, false
		}

		http.SetCookie(*w, &http.Cookie{Name: "x-idp-profile", Value: encoded})

		return profile, ok
	}

	return Profile{}, false
}
