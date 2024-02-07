package config

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatal("load dotenv failed: ", err)
	}

	err = godotenv.Load(path)
	if err != nil {
		log.Fatal("load dotenv failed: ", err)
	}

	return &config{
		app: &app{
			host:         envMap["APP_HOST"],
			port:         loadToInt(envMap, "APP_PORT"),
			name:         envMap["APP_NAME"],
			version:      envMap["APP_VERSION"],
			readTimeout:  loadToSecond(envMap, "APP_READ_TIMEOUT"),
			writeTimeout: loadToSecond(envMap, "APP_WRITE_TIMEOUT"),
			bodyLimit:    loadToInt(envMap, "APP_BODY_LIMIT"),
			fileLimit:    loadToInt(envMap, "APP_FILE_LIMIT"),
			gcpBucket:    envMap["gcpBucket"],
		},
		db: &db{
			host:           envMap["DB_HOST"],
			port:           loadToInt(envMap, "DB_PORT"),
			protocol:       envMap["DB_PROTOCOL"],
			username:       envMap["DB_USERNAME"],
			password:       envMap["DB_PASSWORD"],
			database:       envMap["DB_DATABASE"],
			sslMode:        envMap["DB_SSL_MODE"],
			maxConnections: loadToInt(envMap, "DB_MAX_CONNECTIONS"),
		},
		jwt: &jwt{
			adminKey:       []byte(envMap["JWT_ADMIN_KEY"]),
			secretKey:      []byte(envMap["JWT_SECRET_KEY"]),
			accessExpires:  loadToSecond(envMap, "JWT_ACCESS_EXPIRES"),
			refreshExpires: loadToSecond(envMap, "JWT_REFRESH_EXPIRES"),
		},
		oauth: &oauth{
			sessionSecret: envMap["SESSION_SECRET"],
			googleKey:     envMap["OAUTH_GOOGLE_KEY"],
			googleSecret:  envMap["OAUTH_GOOGLE_SECRET"],
			githubKey:     envMap["OAUTH_GITHUB_KEY"],
			githubSecret:  envMap["OAUTH_GITHUB_SECRET"],
		},
	}
}

func loadToInt(envMap map[string]string, key string) int {
	val, err := strconv.Atoi(envMap[key])
	if err != nil {
		log.Fatalf("load %s failed: %s", key, err)
	}

	return val
}

func loadToSecond(envMap map[string]string, key string) time.Duration {
	val, err := strconv.Atoi(envMap[key])
	if err != nil {
		log.Fatalf("load %s failed: %s", key, err)
	}

	return time.Duration(val) * time.Second
}

// ------------------------------------------------------------- //

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
	OAuth() IOAuthConfig
}

type config struct {
	app   *app
	db    *db
	jwt   *jwt
	oauth *oauth
}

func (c *config) App() IAppConfig     { return c.app }
func (c *config) Db() IDbConfig       { return c.db }
func (c *config) Jwt() IJwtConfig     { return c.jwt }
func (c *config) OAuth() IOAuthConfig { return c.oauth }

// ------------------------------------------------------------- //

type IAppConfig interface {
	Url() string
	Host() string
	Port() int
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
	GCPBucket() string
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int // bytes
	fileLimit    int // bytes
	gcpBucket    string
}

func (a *app) Url() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Host() string                { return a.host }
func (a *app) Port() int                   { return a.port }
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) FileLimit() int              { return a.fileLimit }
func (a *app) GCPBucket() string           { return a.gcpBucket }

// ------------------------------------------------------------- //

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
}

type db struct {
	host           string
	port           int
	protocol       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

func (d *db) Url() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.host, d.port, d.username, d.password, d.database, d.sslMode,
	)
}
func (d *db) MaxOpenConns() int { return d.maxConnections }

// ------------------------------------------------------------- //

type IJwtConfig interface {
	AdminKey() []byte
	SecretKey() []byte
	AccessExpires() time.Duration
	RefreshExpires() time.Duration
}

type jwt struct {
	adminKey       []byte
	secretKey      []byte
	accessExpires  time.Duration
	refreshExpires time.Duration
}

func (j *jwt) AdminKey() []byte              { return j.adminKey }
func (j *jwt) SecretKey() []byte             { return j.secretKey }
func (j *jwt) AccessExpires() time.Duration  { return j.accessExpires }
func (j *jwt) RefreshExpires() time.Duration { return j.refreshExpires }

// ------------------------------------------------------------- //

type IOAuthConfig interface {
	SessionSecret() string
	GoogleKey() string
	GoogleSecret() string
	GithubKey() string
	GithubSecret() string
}

type oauth struct {
	sessionSecret string
	googleKey     string
	googleSecret  string
	githubKey     string
	githubSecret  string
}

func (o *oauth) SessionSecret() string { return o.sessionSecret }
func (o *oauth) GoogleKey() string     { return o.googleKey }
func (o *oauth) GoogleSecret() string  { return o.googleSecret }
func (o *oauth) GithubKey() string     { return o.githubKey }
func (o *oauth) GithubSecret() string  { return o.githubSecret }
