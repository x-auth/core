package authenticators

import (
	"strings"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

func getConfig(authenticator string, realmId string) map[string]string {
	for _, realm := range helpers.Config.Realms {
		if realm.Name == realmId {
			for _, auth := range helpers.Config.Authenticators {
				if auth.Type == authenticator {
					return helpers.ReduceConfig(auth.Config, realm.Config)
				}
			}
		}
	}
	return nil
}

func Login(identifier string, password string, preflightRealm string) (Profile, bool) {
	// check if a valid split char is in the identifier
	ok, splitChar := helpers.SliceContains(identifier, helpers.Config.SplitCharacters)
	if !ok {
		logger.Error.Println("No valid split character in identifier")
		return Profile{}, false
	}

	// split the identifier in username and realm
	idSlice := strings.Split(identifier, splitChar)
	username := identifier
	realmName := idSlice[1]

	// get the realm and autheticator
	var realmObj helpers.Realm
	for _, realm := range helpers.Config.Realms {
		if realm.Name == preflightRealm {
			realmObj = realm
		}
	}

	// quit if the input is wrong
	if realmName != realmObj.Identifier {
		logger.Error.Println("realm did not match preflight: " + realmName + " " + realmObj.Identifier)
		return Profile{}, false
	}

	var authenticator string
	for _, auth := range helpers.Config.Authenticators {
		if auth.Name == realmObj.Authenticator {
			authenticator = auth.Type
		}
	}

	// authenticate using the right authenticator
	if authenticator == "mock" {
		return mock(username, password, getConfig(authenticator, preflightRealm))
	} else if authenticator == "ldap" {
		profile, ok := ldap(username, password, getConfig(authenticator, preflightRealm))
		if !ok {
			logger.Warning.Println("Login failed, Username or password wrong")
			return Profile{}, false
		}

		return profile, ok
	}

	return Profile{}, false
}
