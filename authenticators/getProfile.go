package authenticators

import (
	"net/http"
	"x-net.at/idp/authenticators/ldap"
	"x-net.at/idp/authenticators/mock"
	"x-net.at/idp/models"
)

func GetProfile(authenticator string, username string, cookie *http.Cookie) models.Profile {
	if authenticator == "mock" {
		return mock.GetProfile(username)
	} else if authenticator == "ldap" {
		return ldap.GetProfile(cookie)
	}
	return models.Profile{}
}
