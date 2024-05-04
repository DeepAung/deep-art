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
	Db    *DbConfig
	Jwt   *JwtConfig
	OAuth *OAuthConfig
}

func (c *Config) Print() {
	fmt.Println("===========================================")

	fmt.Println("App")
	fmt.Println("- Address: ", c.App.Address)
	fmt.Println("- Timeout: ", c.App.Timeout)
	fmt.Println("- BodyLimit: ", c.App.BodyLimit)
	fmt.Println("- FileLimit: ", c.App.FileLimit)
	fmt.Println("- CorsOrigins: ", c.App.CorsOrigins)
	fmt.Println("- GcpBucket: ", c.App.GcpBucket)

	fmt.Println("Db")
	fmt.Println("- Url: ", c.Db.Url)

	fmt.Println("Jwt")
	fmt.Println("- SecretKey: ", c.Jwt.SecretKey)
	fmt.Println("- AccessExpires: ", c.Jwt.AccessExpires)
	fmt.Println("- RefreshExpires: ", c.Jwt.RefreshExpires)

	fmt.Println("OAuth")
	fmt.Println("- GoogleKey: ", c.OAuth.GoogleKey)
	fmt.Println("- GoogleSecret: ", c.OAuth.GoogleSecret)
	fmt.Println("- GithubKey: ", c.OAuth.GithubKey)
	fmt.Println("- GithubSecret: ", c.OAuth.GithubSecret)

	fmt.Println("===========================================")
}

type AppConfig struct {
	Address     string
	Timeout     time.Duration
	BodyLimit   string
	FileLimit   string
	CorsOrigins []string
	GcpBucket   string
}

type DbConfig struct {
	Url string
}

type JwtConfig struct {
	SecretKey      string
	AccessExpires  time.Duration
	RefreshExpires time.Duration
}

type OAuthConfig struct {
	GoogleKey    string
	GoogleSecret string
	GithubKey    string
	GithubSecret string
}

func loadEnvPath() string {
	if len(os.Args) == 1 {
		return ""
	}
	return os.Args[1]
}

func NewConfig() *Config {
	path := loadEnvPath()
	if path != "" {
		err := godotenv.Load(path)
		if err != nil {
			log.Fatal("config.go: load env file failed: ", err.Error())
		}
	}

	return &Config{
		App: &AppConfig{
			Address:     os.Getenv("APP_ADDRESS"),
			Timeout:     getAsDuration("APP_TIMEOUT"),
			BodyLimit:   os.Getenv("APP_BODY_LIMIT"),
			FileLimit:   os.Getenv("APP_FILE_LIMIT"),
			CorsOrigins: strings.Split(os.Getenv("APP_CORS_ORIGINS"), " "),
			GcpBucket:   os.Getenv("APP_GCP_BUCKET"),
		},
		Db: &DbConfig{
			Url: os.Getenv("DB_URL"),
		},
		Jwt: &JwtConfig{
			SecretKey:      os.Getenv("JWT_SECRET_KEY"),
			AccessExpires:  getAsDuration("JWT_ACCESS_EXPIRES"),
			RefreshExpires: getAsDuration("JWT_REFRESH_EXPIRES"),
		},
		OAuth: &OAuthConfig{
			GoogleKey:    os.Getenv("OAUTH_GOOGLE_KEY"),
			GoogleSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
			GithubKey:    os.Getenv("OAUTH_GITHUB_KEY"),
			GithubSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
		},
	}
}

func getAsDuration(key string) time.Duration {
	val := os.Getenv(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("config.go: convert string to int error. (\"%s\"=\"%s\")\n", key, val)
	}

	return time.Duration(num) * time.Second
}
