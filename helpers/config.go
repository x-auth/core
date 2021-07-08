package helpers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type realm struct {
	Identifier    string            `yaml:"identifier"`
	Authenticator string            `yaml:"authenticator"`
	Default       bool              `yaml:"default"`
	SkipConsent   bool              `yaml:"skip_consent"`
	Config        map[string]string `yaml:"config,omitempty"`
}

type authenticator struct {
	Type   string            `yaml:"type"`
	Name   string            `yaml:"name"`
	Config map[string]string `yaml:"config"`
}

type conf struct {
	Debug           bool
	Host            string          `yaml:"host"`
	HydraURL        string          `yaml:"hydra_url"`
	KratosURL       string          `yaml:"kratos_url"`
	RememberFor     int64           `yaml:"remember_for"`
	SplitCharacters []string        `yaml:"split_characters"`
	Authenticators  []authenticator `yaml:"authenticators"`
	Realms          []realm         `yaml:"realms"`
}

var Config conf

func LoadConfig() {
	// load the config
	yamlFile, err := ioutil.ReadFile("/etc/idp/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal yaml
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatal(err)
	}

	//set debug mode
	Config.Debug = true

	// validate split characters
	punctuations := ".,:;-!?"
	for _, splitChar := range Config.SplitCharacters {
		if strings.Contains(punctuations, splitChar) {
			log.Fatal("Error: " + splitChar + " is not allowed as a split character")
		}
	}
}
