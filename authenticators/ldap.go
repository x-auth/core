package authenticators

import (
	"crypto/tls"
	"fmt"
	ldap3 "github.com/go-ldap/ldap/v3"
	"net/http"
	"strconv"
	"strings"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

// helper to get ldap attribute by name
func getAttr(attrs []*ldap3.EntryAttribute, name string) []string {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr.Values
		}
	}
	return nil
}

// helper to get all group cns as strings
func getGroups(groupCNs []string) []string {
	var groups []string
	for _, groupCN := range groupCNs {
		groups = append(groups, strings.Split(strings.Split(groupCN, ",")[0], "=")[1])
	}
	return groups
}

func ldap(username string, password string, config map[string]string) (Profile, bool) {
	useTLS, err := strconv.ParseBool(config["use_tls"])
	if err != nil {
		logger.Error.Println(err.Error())
		return Profile{}, false
	}

	// check if ssl/tls cert verification should be skipped
	var conn *ldap3.Conn
	skipVerify, err := strconv.ParseBool(config["skip_verify"])
	if err != nil {
		logger.Error.Println(err.Error())
		return Profile{}, false
	}

	// check if the ldap-server uses tls
	if useTLS {
		// Dial without encryption
		conn, err = ldap3.DialURL("ldap://" + config["host"])
		if err != nil {
			logger.Error.Println(err.Error())
			return Profile{}, false
		}

		// upgrade connection to TLS
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: skipVerify})
		if err != nil {
			logger.Error.Println(err.Error())
			return Profile{}, false
		}
	} else {
		// configure and setup ssl
		tlsConf := &tls.Config{InsecureSkipVerify: skipVerify}
		conn, err = ldap3.DialTLS("tcp", config["host"], tlsConf)

		if err != nil {
			logger.Error.Println(err.Error())
			return Profile{}, false
		}
	}

	// bind with the bind dn
	err = conn.Bind(config["bind_dn"], config["bind_pw"])
	if err != nil {
		logger.Error.Println(err.Error())
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
		logger.Error.Println(err.Error())
		return Profile{}, false
	}

	// get the first Entry of the searchResult
	if len(searchResult.Entries) != 1 {
		logger.Error.Println("User does not exist or too many entries returned")
		return Profile{}, false
	}
	userdn := searchResult.Entries[0].DN

	// Bind as the user to verify password
	err = conn.Bind(userdn, password)
	if err != nil {
		logger.Error.Println(err.Error())
		return Profile{}, false
	}

	// search for the user dn with the users bind, to get the profile information
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
		logger.Error.Println(err.Error())
		return Profile{}, false
	}

	// get the first Entry of the searchResult
	if len(userSearchResult.Entries) != 1 {
		logger.Error.Println("User does not exist or too many entries returned")
		return Profile{}, false
	}

	// parse the ldap entry to the internal Profile struct
	userAttrs := userSearchResult.Entries[0].Attributes
	groups := getGroups(getAttr(userAttrs, config["groups"]))
	profile := Profile{Name: getAttr(userAttrs, config["name"])[0], Email: getAttr(userAttrs, config["email"])[0], Groups: groups}

	// auth successful, profile cookie is set in the http handler
	return profile, true
}

func getLdapProfile(cookie *http.Cookie) Profile {
	// get the profile information from the x-idp-profile cookie
	value := new(Profile)
	err := helpers.SecureCookie.Decode("x-idp-profile", cookie.Value, &value)
	if err != nil {
		logger.Error.Println(err)
		return Profile{}
	}

	return *value
}
