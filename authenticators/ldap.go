package authenticators

import (
	"crypto/tls"
	"fmt"
	ldap3 "github.com/go-ldap/ldap/v3"
	"log"
	"net/http"
	"strconv"
	"x-net.at/idp/helpers"
)

func getAttr(attrs []*ldap3.EntryAttribute, name string) []string {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr.Values
		}
	}
	return nil
}

func ldap(username string, password string, config map[string]string) (Profile, bool) {
	useTLS, err := strconv.ParseBool(config["use_tls"])
	if err != nil {
		return Profile{}, false
	}

	var conn *ldap3.Conn
	skipVerify, err := strconv.ParseBool(config["skip_verify"])
	if err != nil {
		log.Println("LDAP Error: " + err.Error())
		return Profile{}, false
	}

	if useTLS {
		conn, err = ldap3.DialURL("ldap://" + config["host"])
		if err != nil {
			log.Println("LDAP Connection Error: " + err.Error())
			return Profile{}, false
		}

		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: skipVerify})
		if err != nil {
			log.Println("LDAP Error: " + err.Error())
			return Profile{}, false
		}
	} else {
		tlsConf := &tls.Config{InsecureSkipVerify: skipVerify}
		conn, err = ldap3.DialTLS("tcp", config["host"], tlsConf)

		if err != nil {
			log.Println("LDAP Error: " + err.Error())
			return Profile{}, false
		}
	}

	// bind with the bind dn
	err = conn.Bind(config["bind_dn"], config["bind_pw"])
	if err != nil {
		log.Println("LDAP Error: " + err.Error())
		return Profile{}, false
	}

	// set up the search for the given username
	searchRequest := ldap3.NewSearchRequest(
		config["base_dn"],
		ldap3.ScopeWholeSubtree, ldap3.NeverDerefAliases,
		0, 15,
		false,
		fmt.Sprintf("(&%s(%s=%s))", config["filter"], config["email"], username),
		[]string{"dn"},
		nil,
	)

	// run the search
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		log.Println("LDAP Error: " + err.Error())
		return Profile{}, false
	}

	// get the first Entry of the searchResult
	if len(searchResult.Entries) != 1 {
		log.Println("LDAP Error: User does not exist or too many entries returned")
		return Profile{}, false
	}
	userdn := searchResult.Entries[0].DN

	// Bind as the user to verify password
	err = conn.Bind(userdn, password)
	if err != nil {
		log.Println("LDAP Error: " + err.Error())
		return Profile{}, false
	}

	// search with the users bind, to get the profile information
	userSearchRequest := ldap3.NewSearchRequest(
		userdn,
		ldap3.ScopeBaseObject,
		ldap3.NeverDerefAliases,
		0, 15,
		false,
		"(dn=*)",
		[]string{config["email"], config["name"], config["groups"]},
		nil,
	)

	// run the search
	userSearchResult, err := conn.Search(userSearchRequest)
	if err != nil {
		log.Println("LDAP Error: " + err.Error())
		return Profile{}, false
	}

	// get the first Entry of the searchResult
	if len(userSearchResult.Entries) != 1 {
		log.Println("LDAP Rebind Error: User does not exist or too many entries returned")
		return Profile{}, false
	}

	userAttrs := userSearchResult.Entries[0].Attributes
	profile := Profile{Name: getAttr(userAttrs, config["name"])[0], Email: getAttr(userAttrs, config["email"])[0], Groups: getAttr(userAttrs, config["groups"])}

	// auth successful
	return profile, true
}

func getLdapProfile(cookie *http.Cookie) Profile {
	value := new(Profile)
	err := helpers.SecureCookie.Decode("x-idp-profile", cookie.Value, &value)
	if err != nil {
		log.Println(err)
		return Profile{}
	}

	return *value
}
