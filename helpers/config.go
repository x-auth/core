/*
 * Copyright (c) 2021 X-Net Services GmbH
 * Info: https://x-net.at
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package helpers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type Realm struct {
	Name          string `yaml:"name"`
	Identifier    string `yaml:"identifier"`
	Authenticator string `yaml:"authenticator"`
	Default       bool   `yaml:"default"`
	SkipConsent   bool   `yaml:"skip_consent"`
}

type conf struct {
	Debug           bool
	Host            string   `yaml:"host"`
	BasePath        string   `yaml:"base_path"`
	HydraURL        string   `yaml:"hydra_url"`
	KratosURL       string   `yaml:"kratos_url"`
	RememberFor     int64    `yaml:"remember_for"`
	Logger          string   `yaml:"logger"`
	SplitCharacters []string `yaml:"split_characters"`
	Realms          []Realm  `yaml:"realms"`
	RedisAddr       string   `yaml:"redis_addr"`
	RedisPassword   string   `yaml:"redis_password"`
}

var Config conf

func LoadConfig() conf {
	// load the config
	yamlFile, err := ioutil.ReadFile("/etc/x-auth/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal yaml
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatal(err)
	}

	//set debug mode and version
	Config.Debug = false

	// validate split characters
	punctuations := ".,:;-!?"
	for _, splitChar := range Config.SplitCharacters {
		if strings.Contains(punctuations, splitChar) {
			log.Fatal("Error: " + splitChar + " is not allowed as a split character")
		}
	}

	return Config
}
