package services

import (
	"github.com/DeepAung/deep-art/api/types"
	"github.com/markbates/goth"
)

func (s *UsersSvc) OAuthSignup(gothUser goth.User) (types.User, error) {
	return types.User{}, nil
	// var req := typestypecs
}
