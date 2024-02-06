package mytoken

import (
	"fmt"
	"time"

	"github.com/DeepAung/deep-art/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType struct {
	Subject   string
	Audience  jwt.ClaimStrings
	ExpiresAt func(cfg config.IJwtConfig) time.Time
	Key       func(cfg config.IJwtConfig) []byte
}

var Access = TokenType{
	Subject:  "access-token",
	Audience: []string{"admin", "user"},
	ExpiresAt: func(cfg config.IJwtConfig) time.Time {
		return time.Now().Add(cfg.AccessExpires())
	},
	Key: func(cfg config.IJwtConfig) []byte {
		return cfg.SecretKey()
	},
}

var Refresh = TokenType{
	Subject:  "refresh-token",
	Audience: []string{"admin", "user"},
	ExpiresAt: func(cfg config.IJwtConfig) time.Time {
		return time.Now().Add(cfg.RefreshExpires())
	},
	Key: func(cfg config.IJwtConfig) []byte {
		return cfg.SecretKey()
	},
}

var ApiKey = TokenType{
	Subject:  "api-key",
	Audience: []string{"admin", "user"},
	ExpiresAt: func(cfg config.IJwtConfig) time.Time {
		return time.Now().AddDate(2, 0, 0)
	},
	Key: func(cfg config.IJwtConfig) []byte {
		return cfg.ApiKey()
	},
}

var Admin = TokenType{
	Subject:  "admin-token",
	Audience: []string{"admin"},
	ExpiresAt: func(cfg config.IJwtConfig) time.Time {
		return time.Now().Add(5 * time.Minute)
	},
	Key: func(cfg config.IJwtConfig) []byte {
		return cfg.AdminKey()
	},
}

type JwtClaims struct {
	Payload *Payload
	jwt.RegisteredClaims
}

type Payload struct {
	UserId  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
}

func GenerateToken(
	cfg config.IJwtConfig,
	tokenType *TokenType,
	payload *Payload,
) (string, error) {

	claims := &JwtClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "deep-art-api",
			Subject:   tokenType.Subject,
			Audience:  tokenType.Audience,
			ExpiresAt: jwt.NewNumericDate(tokenType.ExpiresAt(cfg)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenType.Key(cfg))
	if err != nil {
		return "", fmt.Errorf("generate token failed: %v", err)
	}

	return tokenString, nil
}

func ParseToken(
	cfg config.IJwtConfig,
	tokenType *TokenType,
	tokenString string,
) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, keyFunc(tokenType.Key(cfg)))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); !ok {
		return nil, fmt.Errorf("invalid claims type")
	} else {
		return claims, nil
	}
}

func VerifyToken(cfg config.IJwtConfig, tokenType *TokenType, tokenString string) error {
	_, err := ParseToken(cfg, tokenType, tokenString)
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
