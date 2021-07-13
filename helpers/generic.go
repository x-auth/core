package helpers

import (
	"strings"
)

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

func ReduceConfig(m1 map[string]string, m2 map[string]string) map[string]string {
	new := make(map[string]string)
	for i1, v1 := range m1 {
		new[i1] = v1
	}

	for i2, v2 := range m2 {
		new[i2] = v2
	}

	return new
}
