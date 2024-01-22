package oauth

import (
	"github.com/DeepAung/deep-art/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func SetupOAuth(cfg config.IOAuthConfig) {
	goth.UseProviders(
		google.New(
			cfg.GoogleKey(),
			cfg.GoogleSecret(),
			"http://127.0.0.1:3000/api/v1/users/google/callback",
		),
		github.New(
			cfg.GithubKey(),
			cfg.GithubSecret(),
			"http://127.0.0.1:3000/api/v1/users/github/callback",
		),
	)
}
