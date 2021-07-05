package helpers

import (
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"log"
)

var SecretKey []byte
var SecureCookie securecookie.SecureCookie

func LoadSecretKey() {
	var err error
	SecretKey, err = ioutil.ReadFile("system/secret.key")
	if err != nil {
		log.Fatal(err)
	}
}

func InitSecureCookie() {
	hashKey := securecookie.GenerateRandomKey(10)
	blockKey := securecookie.GenerateRandomKey(32)

	SecureCookie = *securecookie.New(hashKey, blockKey)
}
