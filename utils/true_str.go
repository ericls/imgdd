package utils

import "strings"

func IsStrTruthy(s string) bool {
	lowerS := strings.ToLower(s)
	return lowerS == "true" || lowerS == "1" || lowerS == "yes" || lowerS == "y" || lowerS == "on"
}
