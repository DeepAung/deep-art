package mytoken

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  = "access-token"
	Refresh = "refresh-token"
)

type JwtClaims struct {
	Payload Payload
	jwt.RegisteredClaims
}

type Payload struct {
	UserId   int
	Username string
}

func GenerateToken(
	tokenType TokenType,
	duration time.Duration,
	secretKey []byte,
	payload Payload,
) (string, error) {

	claims := &JwtClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "deep-art-api",
			Subject:   string(tokenType),
			Audience:  []string{"users", "admin"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("generate token failed: %v", err)
	}

	return tokenString, nil
}

func ParseToken(
	tokenType TokenType,
	secretKey []byte,
	tokenString string,
) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, keyFunc(secretKey))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); !ok {
		return nil, errors.New("invalid claims type")
	} else {
		return claims, nil
	}
}

func VerifyToken(tokenType TokenType, secretKey []byte, tokenString string) error {
	_, err := ParseToken(tokenType, secretKey, tokenString)
	return err
}

func keyFunc(key []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return key, nil
	}
}
