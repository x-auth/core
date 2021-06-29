package authenticators

type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Group     string `json:"group"`
}

func mock(username string, password string) (bool, Profile) {
	if username == "nl@x-net.at" && password == "foobar" {
		return true, Profile{"Nicholas", "Lamprecht", "nl@x-net.at", "Admins"}
	} else {
		return false, Profile{}
	}
}

func getMockProfile(username string) Profile {
	return Profile{"Nicholas", "Lamprecht", "nl@x-net.at", "Admins"}
}
