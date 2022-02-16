package authenticators

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"x-net.at/idp/logger"
)

type authenticator struct {
	Name   string            `yaml:"name"`
	Plugin string            `yaml:"plugin"`
	Config map[string]string `yaml:"config"`
}

type authConfig struct {
	PluginDir      string          `yaml:"plugin_dir"`
	Authenticators []authenticator `yaml:"authenticators"`
}

func LoadAuthConfig() authConfig {
	var cfg authConfig
	yamlFile, err := ioutil.ReadFile("/etc/x-idp/plugins.yaml")
	if err != nil {
		logger.Log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		logger.Log.Fatal(err)
	}

	return cfg
}

func ValidateRealm(authenticatorName string) bool {
	cfg := LoadAuthConfig()
	for _, configured := range cfg.Authenticators {
		if configured.Name == authenticatorName {
			return true
		}
	}

	return false
}
