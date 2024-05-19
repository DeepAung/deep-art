package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
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
