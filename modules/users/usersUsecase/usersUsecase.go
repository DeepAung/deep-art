package usersUsecase

import (
	"fmt"
	"net/http"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/users"
	"github.com/DeepAung/deep-art/modules/users/usersRepository"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

type IUsersUsecase interface {
	Register(req *users.RegisterReq, isAdmin bool) (*users.UserPassport, error)
	Login(req *users.LoginReq) (*users.UserPassport, error)
	Logout(userId, tokenId int) error
	GetUserPassport(user *users.User) (*users.UserPassport, error)
	GetGothUser(c *fiber.Ctx) (*goth.User, error)
	GetUserByOAuth(social users.SocialEnum, socialId string) (bool, *users.User, error)
	CreateOAuth(req *users.OAuthCreateReq) error
	DeleteOAuth(req *users.OAuthReq) error
	RefreshTokens(refreshToken string, userId int) (*users.Token, error)
	GetUserEmailById(userId int) (string, error)
	// Login(req *users.LoginReq) (*users.UserPassport, error) // ???
	// Logout(id int) error
	// RefreshTokens(refreshToken string, id int) (*users.Token, error)
	// GetUserProfile(id int) (*users.User, error)
	// UpdateUserProfile(req *users.UpdateReq) (*users.User, error)
	// DeleteUser(id int) error
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepository.IUsersRepository
}

func NewUsersUsecase(
	cfg config.IConfig,
	usersRepository usersRepository.IUsersRepository,
) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) Register(req *users.RegisterReq, isAdmin bool) (*users.UserPassport, error) {
	if err := req.HashPassword(); err != nil {
		return nil, err
	}

	return u.usersRepository.CreateUser(req, isAdmin)
}

func (u *usersUsecase) Login(req *users.LoginReq) (*users.UserPassport, error) {
	user, err := u.usersRepository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !utils.ComparePassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	return u.GetUserPassport(&users.User{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
		IsAdmin:   user.IsAdmin,
	})
}

// TODO: should we pass accessToken too??? this user might delete tokenId from anoter device login
func (u *usersUsecase) Logout(userId, tokenId int) error {
	return u.usersRepository.DeleteToken(userId, tokenId)
}

func (u *usersUsecase) GetUserPassport(user *users.User) (*users.UserPassport, error) {
	payload := &mytoken.Payload{
		UserId:  user.Id,
		IsAdmin: user.IsAdmin,
	}

	accessToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Access, payload)
	if err != nil {
		return nil, err
	}

	refreshToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Refresh, payload)
	if err != nil {
		return nil, err
	}

	tokenId, err := u.usersRepository.CreateToken(user.Id, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	passport := &users.UserPassport{
		User: &users.User{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
			IsAdmin:   user.IsAdmin,
		},
		Token: &users.Token{
			Id:           tokenId,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	return passport, nil
}

func (u *usersUsecase) GetGothUser(c *fiber.Ctx) (*goth.User, error) {
	var gothUser goth.User
	var err error
	adaptor.HTTPHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		gothUser, err = gothic.CompleteUserAuth(res, req)
	})(c)

	return &gothUser, err
}

func (u *usersUsecase) GetUserByOAuth(
	social users.SocialEnum,
	socialId string,
) (bool, *users.User, error) {
	return u.usersRepository.GetUserByOAuth(social, socialId)
}

func (u *usersUsecase) CreateOAuth(req *users.OAuthCreateReq) error {
	if u.usersRepository.HasOAuth(&users.OAuthReq{UserId: req.UserId, Social: req.Social}) {
		return fmt.Errorf("oauth already connected")
	}

	return u.usersRepository.CreateOAuth(req)
}

func (u *usersUsecase) DeleteOAuth(req *users.OAuthReq) error {
	return u.usersRepository.DeleteOAuth(req)
}

// compare refresh token with database
// verify refresh token???
// gen new one
func (u *usersUsecase) RefreshTokens(
	refreshToken string,
	userId int,
) (*users.Token, error) {
	tokenInfo, err := u.usersRepository.GetTokenInfo(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if tokenInfo.UserId != userId {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// TODO: if user change isAdmin in database, this claims is not up-to-date. this value persist forever...
	claims, err := mytoken.ParseToken(u.cfg.Jwt(), &mytoken.Refresh, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("parse token failed: %v", err)
	}

	newAccessToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Access, claims.Payload)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Refresh, claims.Payload)
	if err != nil {
		return nil, err
	}

	token := &users.Token{
		Id:           tokenInfo.Id,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	if err := u.usersRepository.UpdateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

func (u *usersUsecase) GetUserEmailById(userId int) (string, error) {
	return u.usersRepository.GetUserEmailById(userId)
}
