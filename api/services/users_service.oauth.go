package services

import (
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/markbates/goth"
)

func (s *UsersSvc) OAuthSignup(gothUser goth.User, redirectTo string) (types.User, error) {
	rawPassword := utils.GenRawPassword(16, true, true)
	password, err := utils.Hash(rawPassword)
	if err != nil {
		return types.User{}, err
	}

	req := types.SignUpReq{
		Username:        gothUser.Name,
		Email:           gothUser.Email,
		Password:        password,
		ConfirmPassword: password,
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

	if err := s.usersRepo.CreateOAuthWithDB(ctx, tx, user.Id, gothUser.Provider); err != nil {
		return types.User{}, err
	}

	return user, tx.Commit()
}
