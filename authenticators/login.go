/*
 * Copyright (c) 2021 X-Net Services GmbH
 * Info: https://x-net.at
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package authenticators

import (
	"x-net.at/idp/authenticators/ldap"
	"x-net.at/idp/authenticators/mock"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
	"x-net.at/idp/models"
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

func Login(identifier string, password string, preflightRealm string) (models.Profile, bool) {
	//// check if a valid split char is in the identifier
	//ok, splitChar := helpers.SliceContains(identifier, helpers.Config.SplitCharacters)
	//if !ok {
	//	logger.Log.Error("No valid split character in identifier")
	//	return models.Profile{}, false
	//}
	//
	//// split the identifier in username and realm
	//idSlice := strings.Split(identifier, splitChar)
	//username := identifier
	//realmName := idSlice[1]

	// get the realm and autheticator
	var realmObj helpers.Realm
	for _, realm := range helpers.Config.Realms {
		if realm.Name == preflightRealm {
			realmObj = realm
		}
	}
	logger.Log.Info("realm:", realmObj.Name)

	// quit if the input is wrong
	//if realmName != realmObj.Identifier || realmObj.Default {
	//	logger.Log.Error("realm did not match preflight: " + realmName + " " + realmObj.Identifier)
	//	return models.Profile{}, false
	//}

	var authenticator string
	for _, auth := range helpers.Config.Authenticators {
		if auth.Name == realmObj.Authenticator {
			authenticator = auth.Type
		}
	}

	// authenticate using the right authenticator
	if authenticator == "mock" {
		return mock.Login(identifier, password, getConfig(authenticator, preflightRealm))
	} else if authenticator == "ldap" {
		logger.Log.Debug(getConfig(authenticator, preflightRealm))
		profile, ok := ldap.Login(identifier, password, getConfig(authenticator, preflightRealm))
		if !ok {
			logger.Log.Error("Login failed, Username or password wrong")
			return models.Profile{}, false
		}

		return profile, ok
	}

	return models.Profile{}, false
}
