package usersHandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/users"
	"github.com/DeepAung/deep-art/modules/users/usersUsecase"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

const (
	registerErr      response.TraceId = "users-001"
	loginErr         response.TraceId = "users-002"
	logoutErr        response.TraceId = "users-003"
	refreshTokensErr response.TraceId = "users-006"

	oauthLoginOrRegisterErr response.TraceId = "users-007"
	oauthConnectErr         response.TraceId = "users-007"
	oauthDisconnectErr      response.TraceId = "users-007"

	oauthCallbackErr      response.TraceId = "users-007"
	loginCallbackErr      response.TraceId = "users-007"
	registerCallbackErr   response.TraceId = "users-007"
	connectCallbackErr    response.TraceId = "users-007"
	disconnectCallbackErr response.TraceId = "users-007"
)

type IUsersHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshTokens(c *fiber.Ctx) error

	OAuthLoginOrRegister(c *fiber.Ctx) error
	OAuthConnect(c *fiber.Ctx) error
	OAuthDisconnect(c *fiber.Ctx) error

	OAuthCallback(c *fiber.Ctx) error
	LoginCallback(c *fiber.Ctx, user *users.User) error
	RegisterCallback(
		c *fiber.Ctx,
		gothUser *goth.User,
		user *users.User,
		social users.SocialEnum,
		socialId string,
	) error
	ConnectCallback(c *fiber.Ctx, gothEmail string, social users.SocialEnum, socialId string) error
	DisconnectCallback(c *fiber.Ctx, gothUser *goth.User) error

	// GetUserProfile(c *fiber.Ctx) error
	// UpdateUserProfile(c *fiber.Ctx) error
	// DeleteUser(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecase.IUsersUsecase
}

func NewUsersHandler(cfg config.IConfig, usersUsecase usersUsecase.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) Register(c *fiber.Ctx) error {
	req := new(users.RegisterReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, registerErr, err.Error())
	}

	if !req.IsEmail() {
		return response.Error(c, fiber.StatusBadRequest, registerErr, "invalid email pattern")
	}

	user, err := h.usersUsecase.Register(req)
	if err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return response.Error(c, fiber.StatusBadRequest, registerErr, err.Error())
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return response.Error(c, fiber.StatusBadRequest, registerErr, err.Error())
		default:
			return response.Error(c, fiber.StatusInternalServerError, registerErr, err.Error())
		}
	}

	return response.Success(c, fiber.StatusCreated, user)
}

func (h *usersHandler) Login(c *fiber.Ctx) error {
	req := new(users.LoginReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, loginErr, err.Error())
	}

	passport, err := h.usersUsecase.Login(req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, loginErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, passport)
}

func (h *usersHandler) Logout(c *fiber.Ctx) error {
	// logout oauth first
	adaptor.HTTPHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_ = gothic.Logout(res, req) // TODO: handle error
	})

	req := new(users.LogoutReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, logoutErr, err.Error())
	}

	userId := c.Locals("userId").(int)

	if err := h.usersUsecase.Logout(userId, req.TokenId); err != nil {
		return response.Error(c, fiber.StatusBadRequest, logoutErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, nil)
}

func (h *usersHandler) RefreshTokens(c *fiber.Ctx) error {
	req := new(users.RefreshTokensReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, refreshTokensErr, err.Error())
	}

	userId := c.Locals("userId").(int)

	token, err := h.usersUsecase.RefreshTokens(req, userId)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, refreshTokensErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, token)
}

func (h *usersHandler) OAuthLoginOrRegister(c *fiber.Ctx) error {
	h.setCallbackCookie(c, LoginOrRegister)
	return adaptor.HTTPHandlerFunc(gothic.BeginAuthHandler)(c)
}

func (h *usersHandler) OAuthConnect(c *fiber.Ctx) error {
	h.setCallbackCookie(c, Connect)
	return adaptor.HTTPHandlerFunc(gothic.BeginAuthHandler)(c)
}

func (h *usersHandler) OAuthDisconnect(c *fiber.Ctx) error {
	h.setCallbackCookie(c, Disconnect)
	return adaptor.HTTPHandlerFunc(gothic.BeginAuthHandler)(c)
}

/*
	  if found oauth {
			  login and return passport (user and token)
		} else { // TODO: implement register part
		  if username has been used {
		    append random string after username
		  } else if email has been used {
		    tell user to connect this oauth from normal login(the email way)
		  } else {
		    register with goth user info and return passport(user and nil token)
		  }
		}
*/
func (h *usersHandler) OAuthCallback(c *fiber.Ctx) error {
	gothUser, err := h.usersUsecase.GetGothUser(c)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, oauthCallbackErr, err.Error())
	}

	social := users.SocialEnum(gothUser.Provider)
	socialId := gothUser.UserID

	callback := h.getCallbackCookie(c)
	switch callback {
	case Connect:
		return h.ConnectCallback(c, gothUser.Email, social, socialId)
	case Disconnect:
		return h.DisconnectCallback(c, gothUser)
	}

	// case LoginOrRegister
	found, user, err := h.usersUsecase.GetUserByOAuth(social, socialId)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, oauthCallbackErr, err.Error())
	}

	if found {
		return h.LoginCallback(c, user)
	}

	return h.RegisterCallback(c, gothUser, user, social, socialId)

}

func (h *usersHandler) LoginCallback(c *fiber.Ctx, user *users.User) error {
	passport, err := h.usersUsecase.GetUserPassport(user)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, loginCallbackErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, passport)
}

func (h *usersHandler) RegisterCallback(
	c *fiber.Ctx,
	gothUser *goth.User,
	user *users.User,
	social users.SocialEnum,
	socialId string,
) error {
	passport, err := h.usersUsecase.Register(&users.RegisterReq{
		Username:  gothUser.Name,
		Email:     gothUser.Email,
		Password:  utils.GenRandomPassword(16),
		AvatarUrl: gothUser.AvatarURL,
	})
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, registerCallbackErr, err.Error())
	}

	err = h.usersUsecase.CreateOAuth(&users.OAuthReq{
		UserId:   user.Id,
		Social:   social,
		SocialId: socialId,
	})
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, registerCallbackErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, passport)
}

// get userid, check if email is the same, then link the id
func (h *usersHandler) ConnectCallback(
	c *fiber.Ctx,
	gothEmail string,
	social users.SocialEnum,
	socialId string,
) error {
	userId, err := strconv.Atoi(c.Cookies("userId"))
	if err != nil {
		return response.Error(
			c,
			fiber.StatusInternalServerError,
			connectCallbackErr,
			"jwt auth failed",
		)
	}

	userEmail, err := h.usersUsecase.GetUserEmailById(userId)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, connectCallbackErr, err.Error())
	}

	if userEmail != gothEmail {
		return response.Error(c, fiber.StatusBadRequest, connectCallbackErr, "invalid user id")
	}

	err = h.usersUsecase.CreateOAuth(&users.OAuthReq{
		UserId:   userId,
		Social:   social,
		SocialId: socialId,
	})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, connectCallbackErr, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, nil)
}

// get userid, check if email is the same, then unlink the id
func (h *usersHandler) DisconnectCallback(c *fiber.Ctx, gothUser *goth.User) error {
	return nil
}

// -------------------------------------------------- //

type CallbackEnum string

const (
	LoginOrRegister CallbackEnum = "LoginOrRegister"
	Connect         CallbackEnum = "Connect"
	Disconnect      CallbackEnum = "Disconnect"
)

func (h *usersHandler) setCallbackCookie(c *fiber.Ctx, value CallbackEnum) {
	c.Cookie(&fiber.Cookie{
		Name:     "callback",
		Value:    string(value),
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute), // TODO: is 5 minute OK?
		Secure:   true,
		HTTPOnly: true,
	})
}

func (h *usersHandler) getCallbackCookie(c *fiber.Ctx) CallbackEnum {
	return CallbackEnum(c.Cookies("callback"))
}
