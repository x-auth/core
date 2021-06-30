package helpers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type realm struct {
	Identifier    string `yaml:"identifier"`
	Authenticator string `yaml:"authenticator"`
}

type authenticator struct {
	Type   string            `yaml:"type"`
	Name   string            `yaml:"name"`
	Config map[string]string `yaml:"config"`
}

type conf struct {
	HydraURL        string          `yaml:"hydra_url"`
	RememberFor     int64           `yaml:"remember_for"`
	SplitCharacters []string        `yaml:"split_characters"`
	Authenticators  []authenticator `yaml:"authenticators"`
	Realms          []realm         `yaml:"realms"`
}

var Config conf

func LoadConfig() {
	yamlFile, err := ioutil.ReadFile("/etc/idp/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
