package utils

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes     = "0123456789"
)

func Hash(str string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), 10)
	if err != nil {
		return "", fmt.Errorf("hash password failed: %v", err)
	}

	return string(hashedPassword), nil
}

func Compare(str, hashedStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	return err == nil
}

func GenRawPassword(length int, hasSpecial, hasNumber bool) string {
	password := make([]byte, length)

	charSet := letterBytes
	if hasSpecial {
		charSet += specialBytes
	}
	if hasNumber {
		charSet += numBytes
	}

	for i := range password {
		password[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(password)
}
