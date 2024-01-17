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
	Register(req *users.RegisterReq) (*users.User, error)
	Login(req *users.LoginReq) (*users.UserPassport, error)
	GetGothUser(c *fiber.Ctx) (*goth.User, error)
	GetUserIdByOAuth(social users.SocialEnum, socialId string) (bool, int, error)
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

func (u *usersUsecase) Register(req *users.RegisterReq) (*users.User, error) {
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

	accessToken, err := mytoken.GenerateToken(u.cfg.Jwt(), mytoken.Access, user.Id)
	if err != nil {
		return nil, err // TODO: should this be fmt.Errorf("generate access token failed")???
	}

	refreshToken, err := mytoken.GenerateToken(u.cfg.Jwt(), mytoken.Refresh, user.Id)
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

func (u *usersUsecase) GetUserIdByOAuth(
	social users.SocialEnum,
	socialId string,
) (bool, int, error) {
	return u.usersRepository.GetUserIdByOAuth(social, socialId)
}
