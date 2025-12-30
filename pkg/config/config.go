package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App   *AppConfig
	DB    *DBConfig
	Jwt   *JwtConfig
	OAuth *OAuthConfig
}

func (c *Config) Print() {
	fmt.Println("===========================================")

	fmt.Println("App")
	fmt.Println("- Address: ", c.App.Address)
	fmt.Println("- Timeout: ", c.App.Timeout)
	fmt.Println("- BodyLimit: ", c.App.BodyLimit)
	fmt.Println("- CorsOrigins: ", c.App.CorsOrigins)
	fmt.Println("- GcpBucket: ", c.App.GcpBucket)
	fmt.Println("- BasePath: ", c.App.BasePath)

	fmt.Println("Jwt")
	fmt.Println("- SecretKey: ", c.Jwt.SecretKey)
	fmt.Println("- AccessExpires: ", c.Jwt.AccessExpires)
	fmt.Println("- RefreshExpires: ", c.Jwt.RefreshExpires)

	fmt.Println("OAuth")
	fmt.Println("- GoogleKey: ", c.OAuth.GoogleKey)
	fmt.Println("- GoogleSecret: ", c.OAuth.GoogleSecret)
	fmt.Println("- GithubKey: ", c.OAuth.GithubKey)
	fmt.Println("- GithubSecret: ", c.OAuth.GithubSecret)
	fmt.Println("- SessionSecret: ", c.OAuth.SessionSecret)

	fmt.Println("===========================================")
}

type AppConfig struct {
	Address     string
	Timeout     time.Duration
	BodyLimit   string
	CorsOrigins []string
	GcpBucket   string
	BasePath    string
}

type DBConfig struct {
	Path string
}

type JwtConfig struct {
	SecretKey      []byte
	AccessExpires  time.Duration
	RefreshExpires time.Duration
}

type OAuthConfig struct {
	GoogleKey     string
	GoogleSecret  string
	GithubKey     string
	GithubSecret  string
	SessionSecret string
}

func loadEnvPath() string {
	if len(os.Args) == 1 {
		return ""
	}
	return os.Args[1]
}

func NewConfig(path string) *Config {
	if path == "" {
		path = loadEnvPath()
	}
	if path != "" {
		_ = godotenv.Load(path)
	}

	return &Config{
		App: &AppConfig{
			Address:     os.Getenv("APP_ADDRESS"),
			Timeout:     getAsDuration("APP_TIMEOUT"),
			BodyLimit:   os.Getenv("APP_BODY_LIMIT"),
			CorsOrigins: strings.Split(os.Getenv("APP_CORS_ORIGINS"), " "),
			GcpBucket:   os.Getenv("APP_GCP_BUCKET"),
			BasePath:    os.Getenv("APP_BASE_PATH"),
		},
		DB: &DBConfig{
			Path: os.Getenv("DB_PATH"),
		},
		Jwt: &JwtConfig{
			SecretKey:      []byte(os.Getenv("JWT_SECRET_KEY")),
			AccessExpires:  getAsDuration("JWT_ACCESS_EXPIRES"),
			RefreshExpires: getAsDuration("JWT_REFRESH_EXPIRES"),
		},
		OAuth: &OAuthConfig{
			GoogleKey:     os.Getenv("OAUTH_GOOGLE_KEY"),
			GoogleSecret:  os.Getenv("OAUTH_GOOGLE_SECRET"),
			GithubKey:     os.Getenv("OAUTH_GITHUB_KEY"),
			GithubSecret:  os.Getenv("OAUTH_GITHUB_SECRET"),
			SessionSecret: os.Getenv("SESSION_SECRET"),
		},
	}
}

func getAsDuration(key string) time.Duration {
	var num int = 0
	var err error

	val := os.Getenv(key)
	if val != "" {
		num, err = strconv.Atoi(val)
		if err != nil {
			log.Fatalf("config.go: convert string to int error. (\"%s\"=\"%s\")\n", key, val)
		}
	}

	return time.Duration(num) * time.Second
}
