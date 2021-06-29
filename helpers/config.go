package helpers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	HydraURL string `yaml:"hydra_url"`
	RememberFor int64 `yaml:"remember_for"`
}

var Config conf

func LoadConfig(){
	yamlFile, err := ioutil.ReadFile("/etc/auth/config.yaml")
	if err != nil{
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil{
		log.Fatal(err)
	}
}
