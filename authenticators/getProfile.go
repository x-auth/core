package authenticators

import "net/http"

func GetProfile(authenticator string, username string, cookie *http.Cookie) Profile {
	if authenticator == "mock" {
		return getMockProfile(username)
	} else if authenticator == "ldap" {
		return getLdapProfile(cookie)
	}
	return Profile{}
}
