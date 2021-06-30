package authenticators

import (
	"log"
	"strings"
	"x-net.at/idp/helpers"
)

func Login(identifier string, password string) (string, Profile, bool) {
	/*loggedIn, profile := mock(username, password)
	return "mock", profile, loggedIn*/

	authenticators := helpers.Config.Authenticators
	realms := helpers.Config.Realms

	ok, splitChar := helpers.SliceContains(identifier, helpers.Config.SplitCharacters)
	if !ok {
		return "", Profile{}, false
	}

	idSlice := strings.Split(identifier, splitChar)
	username := idSlice[0]
	realmName := idSlice[1]

	for _, realm := range realms {
		if realm.Identifier == realmName {
			for _, authCfg := range authenticators {
				if authCfg.Name == realm.Authenticator {
					// config := authCfg.Config
					log.Println("using " + authCfg.Type + " authenticator with realm " + realm.Identifier)
					if authCfg.Type == "mock" {
						name, profile, ok := mock(username, password)
						if !ok {
							continue
						} else {
							return name, profile, ok
						}
					}
				}
			}
		}
	}
	return "", Profile{}, false
}
