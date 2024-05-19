package utils

import "regexp"

func IsEmail(email string) bool {
	matched, err := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(email))
	if err != nil {
		return false
	}
	return matched
}
