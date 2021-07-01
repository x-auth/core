package authenticators

import (
	"strings"
	"x-net.at/idp/helpers"
)

func Login(identifier string, password string) (string, Profile, bool) {

	// get the needed config values
	authenticators := helpers.Config.Authenticators
	realms := helpers.Config.Realms

	// check if a valid split char is in the identifier
	ok, splitChar := helpers.SliceContains(identifier, helpers.Config.SplitCharacters)

	// split the identifier in username and realm
	idSlice := strings.Split(identifier, splitChar)
	username := idSlice[0]
	realmName := idSlice[1]

	// default fallback
	if !ok {
		for _, realm := range realms {
			if realm.Default {
				for _, authCfg := range authenticators {
					if authCfg.Name == realm.Authenticator {
						if authCfg.Type == "mock" {
							cfg := helpers.ReduceConfig(authCfg.Config, realm.Config)
							name, profile, ok := mock(username, password, cfg)
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
	}

	// iterate over all configured realms
	for _, realm := range realms {
		// check if the realm in the identifier matches a configured realm
		if realm.Identifier == realmName {
			for _, authCfg := range authenticators {
				if authCfg.Name == realm.Authenticator {
					if authCfg.Type == "mock" {
						cfg := helpers.ReduceConfig(authCfg.Config, realm.Config)
						name, profile, ok := mock(username, password, cfg)
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
