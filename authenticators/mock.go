package authenticators

func mock(username string, password string, config map[string]string) (Profile, bool) {
	if username == config["username"] && password == config["password"] {
		return Profile{"Nicholas Lamprecht", "nl@x-net.at", []string{"Admins"}}, true
	} else {
		return Profile{}, false
	}
}

func getMockProfile(username string) Profile {
	return Profile{"Nicholas Lamprecht", "nl@x-net.at", []string{"Admins"}}
}
