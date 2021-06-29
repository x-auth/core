package authenticators

func GetProfile(authenticator string, username string)Profile{
	if authenticator == "mock" {
		return getMockProfile(username)
	}
	return Profile{}
}
