package authenticators

type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Group     string `json:"group"`
}

func mock(username string, password string, config map[string]string) (string, Profile, bool) {
	if username == config["username"] && password == config["password"] {
		return "mock", Profile{"Nicholas", "Lamprecht", "nl@x-net.at", "Admins"}, true
	} else {
		return "mock", Profile{}, false
	}
}

func getMockProfile(username string) Profile {
	return Profile{"Nicholas", "Lamprecht", "nl@x-net.at", "Admins"}
}
