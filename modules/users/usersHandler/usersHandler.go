package usersHandler

import (
	"net/http"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/users"
	"github.com/DeepAung/deep-art/modules/users/usersUsecase"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/markbates/goth/gothic"
)

const (
	registerErr          response.TraceId = "users-001"
	loginErr             response.TraceId = "users-002"
	logoutErr            response.TraceId = "users-003"
	authenticateOAuthErr response.TraceId = "users-004"
	callbackOAuthErr     response.TraceId = "users-005"
)

type IUsersHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error

	AuthenticateOAuth(c *fiber.Ctx) error
	CallbackOAuth(c *fiber.Ctx) error
	// RefreshTokens(c *fiber.Ctx) error
	// ConnectOAuth(c *fiber.Ctx) error
	// DisconnectOAuth(c *fiber.Ctx) error
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

func (h *usersHandler) AuthenticateOAuth(c *fiber.Ctx) error {
	return adaptor.HTTPHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})(c)
}

/*
	if found oauth {
	  login and return passport (user and token)
	} else { // register part

	  if username has been used {
	    append random string after username
	  } else if email has been used {
	    tell user to connect this oauth from normal login(the email way)
	  } else {
	    register with goth user info and return passport(user and nil token)
	  }
	}
*/
func (h *usersHandler) CallbackOAuth(c *fiber.Ctx) error {
	gothUser, err := h.usersUsecase.GetGothUser(c)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, callbackOAuthErr, err.Error())
	}

	social := users.SocialEnum(gothUser.Provider)
	socialId := gothUser.UserID

	found, user, err := h.usersUsecase.GetUserByOAuth(social, socialId)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, callbackOAuthErr, err.Error())
	}

	if found {
		passport, err := h.usersUsecase.GetUserPassport(user)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, callbackOAuthErr, err.Error())
		}

		return response.Success(c, fiber.StatusOK, passport)
	}

	passport, err := h.usersUsecase.Register(&users.RegisterReq{
		Username:  gothUser.Name,
		Email:     gothUser.Email,
		Password:  utils.GenRandomPassword(16),
		AvatarUrl: gothUser.AvatarURL,
	})
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, callbackOAuthErr, err.Error())
	}

	err = h.usersUsecase.CreateOAuth(&users.OAuthReq{
		UserId:   user.Id,
		Social:   social,
		SocialId: socialId,
	})
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, callbackOAuthErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, passport)
}
