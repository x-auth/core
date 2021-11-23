package ldap

import (
	ldap3 "github.com/go-ldap/ldap/v3"
	"strings"
)

func getAttr(attrs []*ldap3.EntryAttribute, name string) []string {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr.Values
		}
	}
	return []string{""}
}

// helper to get all group cns as strings
func getGroups(groupCNs []string) []string {
	var groups []string
	for _, groupCN := range groupCNs {
		parsedArray := strings.Split(strings.Split(groupCN, ",")[0], "=")

		if len(parsedArray) == 0 {
			continue
		}

		groups = append(groups, parsedArray[1])
	}
	return groups
}
