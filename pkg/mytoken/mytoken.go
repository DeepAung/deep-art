package mytoken

import (
	"fmt"
	"time"

	"github.com/DeepAung/deep-art/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access-token"
	Refresh TokenType = "refresh-token"
)

type JwtClaims struct {
	Payload *Payload
	jwt.RegisteredClaims
}

type Payload struct {
	Id int `json:"id"`
}

func GenerateToken(
	cfg config.IJwtConfig,
	tokenType TokenType,
	userId int,
) (string, error) {

	expDuration := getExpDuration(cfg, tokenType)
	claims := &JwtClaims{
		Payload: &Payload{
			Id: userId,
		},
		RegisteredClaims: jwt.RegisteredClaims{ // TODO: should we have this field or ignore it???
			Issuer:    "deep-art-claims",
			Subject:   string(tokenType),
			Audience:  []string{"user"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(cfg.SecretKey())
	if err != nil {
		return "", fmt.Errorf("sign token to string failed: %v", err)
	}

	return tokenString, nil
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, returnKeyFunc(cfg))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); !ok {
		return nil, fmt.Errorf("invalid claims type")
	} else {
		return claims, nil
	}
}

func returnKeyFunc(cfg config.IJwtConfig) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return cfg.SecretKey(), nil
	}
}

func getExpDuration(cfg config.IJwtConfig, tokenType TokenType) time.Duration {
	switch tokenType {
	case Access:
		return cfg.AccessExpires()
	case Refresh:
		return cfg.RefreshExpires()
	default:
		return cfg.AccessExpires() // TODO: maybe handle some error
	}
}
