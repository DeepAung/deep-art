package usersHandler

import (
	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/users"
	"github.com/DeepAung/deep-art/modules/users/usersUsecase"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/markbates/goth/gothic"
)

const (
	registerErr          response.TraceId = "users-001"
	loginErr             response.TraceId = "users-002"
	authenticateOAuthErr response.TraceId = "users-003"
	callbackOAuthErr     response.TraceId = "users-004"
)

type IUsersHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	AuthenticateOAuth(c *fiber.Ctx) error
	CallbackOAuth(c *fiber.Ctx) error

	// Logout(c *fiber.Ctx) error
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

func (h *usersHandler) AuthenticateOAuth(c *fiber.Ctx) error {
	return adaptor.HTTPHandlerFunc(gothic.BeginAuthHandler)(c)
}

/*
1) return username email profile for register handler
2) return token for login
*/
func (h *usersHandler) CallbackOAuth(c *fiber.Ctx) error {
	_, err := h.usersUsecase.GetGothUser(c)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, callbackOAuthErr, err.Error())
	}

	return nil

	/*
		found, userId, err := h.usersUsecase.GetUserIdByOAuth(
			users.SocialEnum(gothUser.Provider),
			gothUser.UserID,
		)
		if err != nil {
			return nil // TODO:
		}

		// TODO:
		if found {
			// login and return passport
			// h.usersUsecase.GetPassport(userId)
			return nil
		} else {
			// register with default values (username, email, profile)
			h.usersUsecase.CreateUser(&users.RegisterReq{
				Username:  gothUser.Name,
				Email:     gothUser.Email,
				Password:  genpass, // TODO: genpassword.GenRandomPassword and Hash it
				AvatarUrl: gothUser.AvatarURL,
			})
		}
	*/
}
