package services

import (
	"errors"
	"net/http"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/markbates/goth"
)

var ErrOAuthSignin = httperror.New(
	"you didn't have OAuth connected. you can either Sign-up with OAuth or Sign-in to connect OAuth",
	http.StatusBadRequest,
)

func (s *UsersSvc) OAuthSignup(gothUser goth.User, redirectTo string) (types.User, error) {
	req := types.SignUpReq{
		Username:        gothUser.Name,
		Email:           gothUser.Email,
		Password:        "",
		ConfirmPassword: "",
		AvatarUrl:       gothUser.AvatarURL,
		RedirectTo:      redirectTo,
	}

	ctx, cancel, tx, err := s.usersRepo.BeginTx()
	defer cancel()
	if err != nil {
		return types.User{}, err
	}

	user, err := s.usersRepo.CreateUserWithDB(ctx, tx, req)
	if err != nil {
		return types.User{}, err
	}

	if err := s.usersRepo.CreateOAuthWithDB(ctx, tx, user.Id, gothUser.Provider, gothUser.UserID); err != nil {
		return types.User{}, err
	}

	return user, tx.Commit()
}

func (s *UsersSvc) OAuthSignin(gothUser goth.User, redirectTo string) (types.Passport, error) {
	has, err := s.usersRepo.HasOAuth(gothUser.UserID, gothUser.Provider)
	if err != nil {
		return types.Passport{}, err
	}
	if !has {
		return types.Passport{}, ErrOAuthSignin
	}

	user, err := s.usersRepo.FindOneUserByEmail(gothUser.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return types.Passport{}, ErrInvalidEmailOrPassword
		}
		return types.Passport{}, err
	}

	return s.generatePassport(user)
}