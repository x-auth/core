package authenticators

import (
	"math/rand"
	"strconv"
	"x-net.at/idp/logger"
)

func mock(username string, password string, config map[string]string) (Profile, bool) {
	if username == config["username"] && password == config["password"] {
		return Profile{strconv.Itoa(rand.Intn(100)), "Nicholas Lamprecht", "nl@x-net.at", []string{"Admins"}}, true
	} else {
		logger.Warning.Println("Login failed, username or password false")
		return Profile{}, false
	}
}

func getMockProfile(username string) Profile {
	return Profile{strconv.Itoa(rand.Intn(100)), "Nicholas Lamprecht", "nl@x-net.at", []string{"Admins"}}
}
