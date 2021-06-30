package helpers

import "strings"

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceContains(str string, list []string) (bool, string) {
	for _, char := range list {
		if strings.Contains(str, char) {
			return true, char
		}
	}
	return false, ""
}
