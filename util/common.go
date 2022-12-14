package util

import (
	"strings"
)

// ContainInArray check whether or not item in array
func ContainInArray(array []string, item string) bool {
	for _, value := range array {
		if strings.EqualFold(value, item) {
			return true
		}
	}
	return false
}
