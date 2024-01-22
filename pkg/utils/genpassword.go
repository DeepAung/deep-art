package utils

import (
	"math/rand"
)

const (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func GenRandomPassword(length int) string {
	password := make([]rune, length)
	lenAllCharSet := len(allCharSet)

	for i := 0; i < length; i++ {
		randIndex := rand.Intn(lenAllCharSet)
		password[i] = rune(allCharSet[randIndex])
	}

	return string(password)
}
