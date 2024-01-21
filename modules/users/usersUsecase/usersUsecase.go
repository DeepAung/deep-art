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
	Register(req *users.RegisterReq) (*users.UserPassport, error)
	Login(req *users.LoginReq) (*users.UserPassport, error)
	Logout(userId, tokenId int) error
	GetUserPassport(user *users.User) (*users.UserPassport, error)
	GetGothUser(c *fiber.Ctx) (*goth.User, error)
	GetUserByOAuth(social users.SocialEnum, socialId string) (bool, *users.User, error)
	CreateOAuth(req *users.OAuthReq) error
	RefreshTokens(req *users.RefreshTokensReq, userId int) (*users.Token, error)
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

func (u *usersUsecase) Register(req *users.RegisterReq) (*users.UserPassport, error) {
	if err := req.HashPassword(); err != nil {
		return nil, err
	}

	return u.usersRepository.CreateUser(req)
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
	})
}

// TODO: should we pass accessToken too??? this user might delete tokenId from anoter device login
func (u *usersUsecase) Logout(userId, tokenId int) error {
	return u.usersRepository.DeleteToken(userId, tokenId)
}

func (u *usersUsecase) GetUserPassport(user *users.User) (*users.UserPassport, error) {
	accessToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Access, user.Id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Refresh, user.Id)
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

func (u *usersUsecase) CreateOAuth(req *users.OAuthReq) error {
	return u.usersRepository.CreateOAuth(req)
}

func (u *usersUsecase) RefreshTokens(
	req *users.RefreshTokensReq,
	userId int,
) (*users.Token, error) {
	tokenInfo, err := u.usersRepository.GetToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if tokenInfo.UserId != userId {
		return nil, fmt.Errorf("invalid refresh token")
	}

	accessToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Access, userId)
	if err != nil {
		return nil, err
	}

	refreshToken, err := mytoken.GenerateToken(u.cfg.Jwt(), &mytoken.Refresh, userId)
	if err != nil {
		return nil, err
	}

	token := &users.Token{
		Id:           tokenInfo.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := u.usersRepository.UpdateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}
