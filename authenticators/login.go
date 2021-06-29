package authenticators

func Login(username string, password string)(string, Profile, bool){
	loggedIn, profile := mock(username, password)
	return "mock", profile, loggedIn
}
