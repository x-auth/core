package mock

import (
	"math/rand"
	"strconv"
	"x-net.at/idp/logger"
	"x-net.at/idp/models"
)

func Login(username string, password string, config map[string]string) (models.Profile, bool) {
	if username == config["username"] && password == config["password"] {
		return models.Profile{strconv.Itoa(rand.Intn(100)), "Nicholas Lamprecht", "nl@x-net.at", []string{"Admins"}}, true
	} else {
		logger.Warning.Println("Login failed, username or password false")
		return models.Profile{}, false
	}
}

func getMockProfile(username string) models.Profile {
	return models.Profile{strconv.Itoa(rand.Intn(100)), "Nicholas Lamprecht", "nl@x-net.at", []string{"Admins"}}
}
