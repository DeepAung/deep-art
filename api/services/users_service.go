package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/DeepAung/deep-art/pkg/utils"
)

var (
	ErrInvalidEmailOrPassword = httperror.New("invalid email or password", http.StatusBadRequest)
	ErrInvalidRefreshToken    = httperror.New("invalid refresh token", http.StatusBadRequest)
)

type UsersSvc struct {
	usersRepo *repositories.UsersRepo
	storer    storer.Storer
	cfg       *config.Config
}

func NewUsersSvc(
	usersRepo *repositories.UsersRepo,
	storer storer.Storer,
	cfg *config.Config,
) *UsersSvc {
	return &UsersSvc{
		usersRepo: usersRepo,
		storer:    storer,
		cfg:       cfg,
	}
}

func (s *UsersSvc) SignIn(email string, password string) (types.Passport, error) {
	user, err := s.usersRepo.FindOneUserWithPasswordByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return types.Passport{}, ErrInvalidEmailOrPassword
		}
		return types.Passport{}, err
	}

	// User only sign-in with OAuth
	if user.Password == "" {
		return types.Passport{}, ErrInvalidEmailOrPassword
	}

	if !utils.Compare(password, user.Password) {
		return types.Passport{}, ErrInvalidEmailOrPassword
	}

	return s.generatePassport(types.User{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
		IsAdmin:   user.IsAdmin,
		Coin:      user.Coin,
	})
}

func (s *UsersSvc) generatePassport(user types.User) (types.Passport, error) {
	payload := mytoken.Payload{
		UserId:   user.Id,
		Username: user.Username,
	}

	accessToken, err := mytoken.GenerateToken(
		mytoken.Access,
		s.cfg.Jwt.AccessExpires,
		s.cfg.Jwt.SecretKey,
		payload,
	)
	if err != nil {
		return types.Passport{}, err
	}

	refreshToken, err := mytoken.GenerateToken(
		mytoken.Refresh,
		s.cfg.Jwt.RefreshExpires,
		s.cfg.Jwt.SecretKey,
		payload,
	)
	if err != nil {
		return types.Passport{}, err
	}

	tokenId, err := s.usersRepo.CreateToken(user.Id, accessToken, refreshToken)
	if err != nil {
		return types.Passport{}, err
	}

	passport := types.Passport{
		User: user,
		Token: types.Token{
			Id:           tokenId,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	return passport, nil
}

func (s *UsersSvc) SignUp(req types.SignUpReq) (types.User, error) {
	var err error
	if req.Password, err = utils.Hash(req.Password); err != nil {
		return types.User{}, err
	}

	return s.usersRepo.CreateUser(req)
}

func (s *UsersSvc) SignOut(userId int, tokenId int) error {
	return s.usersRepo.DeleteToken(userId, tokenId)
}

func (s *UsersSvc) UpdateTokens(userId int, refreshToken string) (types.Token, error) {
	tokenId, err := s.usersRepo.FindOneTokenId(userId, refreshToken)
	if err != nil {
		if errors.Is(err, repositories.ErrTokenNotFound) {
			return types.Token{}, ErrInvalidRefreshToken
		} else {
			return types.Token{}, err
		}
	}

	claims, err := mytoken.ParseToken(mytoken.Refresh, s.cfg.Jwt.SecretKey, refreshToken)
	if err != nil {
		return types.Token{}, err
	}

	newAccessToken, err := mytoken.GenerateToken(
		mytoken.Access,
		s.cfg.Jwt.AccessExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return types.Token{}, err
	}

	newRefreshToken, err := mytoken.GenerateToken(
		mytoken.Refresh,
		s.cfg.Jwt.RefreshExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return types.Token{}, err
	}

	err = s.usersRepo.UpdateTokens(tokenId, newAccessToken, newRefreshToken)
	if err != nil {
		return types.Token{}, err
	}

	return types.Token{
		Id:           tokenId,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *UsersSvc) HasAccessToken(userId int, accessToken string) (bool, error) {
	return s.usersRepo.HasAccessToken(userId, accessToken)
}

func (s *UsersSvc) HasRefreshToken(userId int, refreshToken string) (bool, error) {
	return s.usersRepo.HasRefreshToken(userId, refreshToken)
}

func (s *UsersSvc) GetUser(id int) (types.User, error) {
	return s.usersRepo.FindOneUserById(id)
}

func (s *UsersSvc) GetCreator(id int) (types.Creator, error) {
	return s.usersRepo.FindOneCreatorById(id)
}

func (s *UsersSvc) UpdateUser(id int, avatar *multipart.FileHeader, req types.UpdateUserReq) error {
	user, err := s.usersRepo.FindOneUserById(id)
	if err != nil {
		return err
	}

	// delete old avatar
	if user.AvatarUrl != "" {
		dest := utils.NewUrlInfoByURL(s.cfg.App.BasePath, user.AvatarUrl).Dest()
		if err := s.storer.DeleteFiles([]string{dest}); err != nil {
			return err
		}
	}

	// upload new avatar
	dir := fmt.Sprint("/users/", id)
	res, err := s.storer.UploadFiles([]*multipart.FileHeader{avatar}, dir)
	if err != nil {
		return err
	}

	// update user field (username, avatarUrl)
	req.AvatarUrl = res[0].Url()
	if err = s.usersRepo.UpdateUser(id, req); err != nil {
		dest := fmt.Sprint("/users/", id, "/", res[0].Filename())
		_ = s.storer.DeleteFiles([]string{dest})
		return err
	}

	return nil
}

func (s *UsersSvc) DeleteUser(id int) error {
	return s.usersRepo.DeleteUser(id)
}

// follower = user
// followee = creator
func (s *UsersSvc) ToggleFollow(followerId, followeeId int) (bool, error) {
	isFollowing, err := s.IsFollowing(followerId, followeeId)
	if err != nil {
		return false, err
	}

	if isFollowing {
		err = s.UnFollow(followerId, followeeId)
	} else {
		err = s.Follow(followerId, followeeId)
	}

	if err != nil {
		return false, err
	}

	return !isFollowing, nil
}

func (s *UsersSvc) IsFollowing(followerId, followeeId int) (bool, error) {
	return s.usersRepo.HasFollow(followerId, followeeId)
}

func (s *UsersSvc) Follow(followerId, followeeId int) error {
	return s.usersRepo.CreateFollow(followerId, followeeId)
}

func (s *UsersSvc) UnFollow(followerId, followeeId int) error {
	return s.usersRepo.DeleteFollow(followerId, followeeId)
}
