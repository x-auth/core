package mock

import (
	"x-net.at/idp/logger"
	"x-net.at/idp/models"
)

func Login(username string, password string, config map[string]string) (models.Profile, bool) {
	if username == config["username"] && password == config["password"] {
		return models.Profile{
			Name:        "Foo Bar",
			FamilyName:  "Bar",
			GivenName:   "Foo",
			NickName:    "foobar",
			Email:       "foobar@example.com",
			PhoneNumber: "000000000",
		}, true
	} else {
		logger.Warning.Println("Login failed, username or password false")
		return models.Profile{}, false
	}
}

func getMockProfile(username string) models.Profile {
	return models.Profile{
		Name:        "Foo Bar",
		FamilyName:  "Bar",
		GivenName:   "Foo",
		NickName:    "foobar",
		Email:       "foobar@example.com",
		PhoneNumber: "000000000",
	}
}
