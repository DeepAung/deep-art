package mytoken

import (
	"fmt"
	"time"

	"github.com/DeepAung/deep-art/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType struct {
	Subject   string
	ExpiresAt func(cfg config.IJwtConfig) time.Time
}

var Access = TokenType{
	Subject: "access-token",
	ExpiresAt: func(cfg config.IJwtConfig) time.Time {
		return time.Now().Add(cfg.AccessExpires())
	},
}

var Refresh = TokenType{
	Subject: "refresh-token",
	ExpiresAt: func(cfg config.IJwtConfig) time.Time {
		return time.Now().Add(cfg.RefreshExpires())
	},
}

type JwtClaims struct {
	Payload *Payload
	jwt.RegisteredClaims
}

type Payload struct {
	Id int `json:"id"`
}

func GenerateToken(
	cfg config.IJwtConfig,
	tokenType *TokenType,
	userId int,
) (string, error) {

	claims := &JwtClaims{
		Payload: &Payload{
			Id: userId,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "deep-art-api",
			Subject:   tokenType.Subject,
			ExpiresAt: jwt.NewNumericDate(tokenType.ExpiresAt(cfg)),
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
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, keyFunc(cfg))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); !ok {
		return nil, fmt.Errorf("invalid claims type")
	} else {
		return claims, nil
	}
}

func keyFunc(cfg config.IJwtConfig) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return cfg.SecretKey(), nil
	}
}
