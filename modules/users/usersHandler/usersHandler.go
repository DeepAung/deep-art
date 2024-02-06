package usersHandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/users"
	"github.com/DeepAung/deep-art/modules/users/usersUsecase"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

/*
case of authentications

1) user register with example@gmail.com
   and try to connect google with example2@gmail.com
   response: this google email is not the same as current email

2) user register with example@gmail.com
   and try to login/register google with example@gmail.com
   response: this email already exist try login and connect it
   or: connect it automatically???

3) user register google with example@gmail.com
   and try to register with example@gmail.com
   reponse: this email is already used
*/

const (
	// users-000, users-001, and so on
	_ = response.TraceId(
		"users-" +
			string('0'+iota/100%10) +
			string('0'+iota/10%10) +
			string('0'+iota/1%10))

	registerErr
	loginErr
	logoutErr
	refreshTokensErr

	oauthLoginOrRegisterErr
	oauthConnectErr
	oauthDisconnectErr

	oauthCallbackErr
	loginCallbackErr
	registerCallbackErr
	connectCallbackErr

	generateAdminTokenErr
)

type CallbackEnum string

const (
	LoginOrRegister CallbackEnum = "LoginOrRegister"
	Connect         CallbackEnum = "Connect"
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
	ConnectCallback(c *fiber.Ctx, email string, social users.SocialEnum, socialId string) error

	GenerateAdminToken(c *fiber.Ctx) error
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
		case "email has been used", "username has been used":
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

	err := mytoken.VerifyToken(h.cfg.Jwt(), &mytoken.Refresh, req.RefreshToken)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, refreshTokensErr, err.Error())
	}

	userId := c.Locals("userId").(int)

	token, err := h.usersUsecase.RefreshTokens(req.RefreshToken, userId)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, refreshTokensErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, token)
}

func (h *usersHandler) OAuthLoginOrRegister(c *fiber.Ctx) error {
	utils.SetCookie(c, "callback", string(LoginOrRegister), 5*time.Minute) // TODO: is 5 minute OK?
	return adaptor.HTTPHandlerFunc(gothic.BeginAuthHandler)(c)
}

func (h *usersHandler) OAuthConnect(c *fiber.Ctx) error {
	utils.SetCookie(c, "callback", string(Connect), 5*time.Minute)
	return adaptor.HTTPHandlerFunc(gothic.BeginAuthHandler)(c)
}

func (h *usersHandler) OAuthDisconnect(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	req := new(users.OAuthDisconnectReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, oauthDisconnectErr, err.Error())
	}

	err := h.usersUsecase.DeleteOAuth(&users.OAuthReq{
		UserId: userId,
		Social: req.Social,
	})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, connectCallbackErr, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, nil)
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

	callback := CallbackEnum(c.Cookies("callback"))
	c.ClearCookie("callback")

	if callback == Connect {
		return h.ConnectCallback(c, gothUser.Email, social, socialId)
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

	err = h.usersUsecase.CreateOAuth(&users.OAuthCreateReq{
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
	email string,
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

	if userEmail != email {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			connectCallbackErr,
			"email in this social connect is not the same as current email",
		)
	}

	err = h.usersUsecase.CreateOAuth(&users.OAuthCreateReq{
		UserId:   userId,
		Social:   social,
		SocialId: socialId,
	})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, connectCallbackErr, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, nil)
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	token, err := mytoken.GenerateToken(h.cfg.Jwt(), &mytoken.Admin, nil)
	if err != nil {
		return response.Error(
			c,
			fiber.StatusInternalServerError,
			generateAdminTokenErr,
			err.Error(),
		)
	}

	return response.Success(c, fiber.StatusCreated, &users.AdminTokenRes{
		AdminToken: token,
	})
}
