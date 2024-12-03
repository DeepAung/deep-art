package services_test

import (
	"database/sql"
	"mime/multipart"
	"testing"
	"time"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/asserts"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/golang-migrate/migrate/v4"
)

var (
	cfg       *config.Config
	testDB    *sql.DB
	migrateDB *migrate.Migrate
	mystorer  storer.Storer
	repo      *repositories.UsersRepo
	svc       *services.UsersSvc
)

func init() {
	testDB, migrateDB = repositories.NewTestDB()
	repositories.ResetDB(migrateDB)

	cfg = config.NewConfig("../../.env.dev")
	mystorer = storer.NewGCPStorer(cfg)
	repo = repositories.NewUsersRepo(testDB, 1*time.Second)
	svc = services.NewUsersSvc(repo, mystorer, cfg)
}

func Test_UserSvc_Signin(t *testing.T) {
	tests := []struct {
		name             string
		email            string
		password         string
		expectedPassport types.Passport
		expectedError    error
	}{
		{
			name:             "no email",
			email:            "",
			password:         "password",
			expectedPassport: types.Passport{},
			expectedError:    services.ErrInvalidEmailOrPassword,
		},
		{
			name:             "no password",
			email:            "i.deepaung@gmail.com",
			password:         "",
			expectedPassport: types.Passport{},
			expectedError:    services.ErrInvalidEmailOrPassword,
		},
		{
			name:             "not found email",
			email:            "invalid-email@gmail.com",
			password:         "",
			expectedPassport: types.Passport{},
			expectedError:    services.ErrInvalidEmailOrPassword,
		},
		{
			name:             "invalid password",
			email:            "i.deepaung@gmail.com",
			password:         "invalid-password",
			expectedPassport: types.Passport{},
			expectedError:    services.ErrInvalidEmailOrPassword,
		},
		{
			name:     "normal signin",
			email:    "i.deepaung@gmail.com",
			password: "password",
			expectedPassport: types.Passport{
				User: types.User{
					Id:        1,
					Username:  "DeepAung",
					Email:     "i.deepaung@gmail.com",
					AvatarUrl: "",
					IsAdmin:   false,
					Coin:      0,
				},
				Token: types.Token{
					Id: 1,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passport, err := svc.SignIn(tt.email, tt.password)

			asserts.EqualError(t, err, tt.expectedError)
			if err == nil {
				asserts.NotEqual(t, "passport.Token.AccessToken", passport.Token.AccessToken, "")
				asserts.NotEqual(t, "passport.Token.RefreshToken", passport.Token.RefreshToken, "")
				passport.Token.AccessToken = ""
				passport.Token.RefreshToken = ""
				asserts.Equal(t, "passport", passport, tt.expectedPassport)
			}
		})
	}
}

func Test_UserSvc_SignOut(t *testing.T) {
	tests := []struct {
		name          string
		userId        int
		tokenId       int
		expectedError error
	}{
		{
			name:          "invalid user id",
			userId:        555,
			tokenId:       1,
			expectedError: services.ErrUserOrTokenNotFound,
		},
		{
			name:          "invalid token id",
			userId:        1,
			tokenId:       555,
			expectedError: services.ErrUserOrTokenNotFound,
		},
		{
			name:          "normal signout",
			userId:        1,
			tokenId:       1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.SignOut(tt.userId, tt.tokenId)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UserSvc_SignUp(t *testing.T) {
	req := types.SignUpReq{
		Username:        "newuser",
		Email:           "newuser@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
		AvatarUrl:       "www.newuser.com/pic",
	}
	expectedUser := types.User{
		Id:        7,
		Username:  "newuser",
		Email:     "newuser@gmail.com",
		AvatarUrl: "www.newuser.com/pic",
		IsAdmin:   false,
		Coin:      0,
	}

	user, err := svc.SignUp(req)

	asserts.EqualError(t, err, nil)

	asserts.Equal(t, "user", user, expectedUser)
}

func Test_UserSvc_GetUser(t *testing.T) {
	_, err := svc.GetUser(1)
	asserts.EqualError(t, err, nil)

	_, err = svc.GetUser(1000)
	asserts.EqualError(t, err, repositories.ErrUserNotFound)
}

func Test_UserSvc_GetCreator(t *testing.T) {
	_, err := svc.GetCreator(1)
	asserts.EqualError(t, err, nil)

	_, err = svc.GetCreator(1000)
	asserts.EqualError(t, err, repositories.ErrUserNotFound)
}

func Test_UserSvc_SetPassword(t *testing.T) {
	err := svc.SetPassword(7, "newpassword")
	asserts.EqualError(t, err, nil)

	err = svc.SetPassword(7, "")
	asserts.EqualError(t, err, services.ErrCannotSetEmptyPassword)
}

func Test_UserSvc_HasPassword(t *testing.T) {
	has, err := svc.HasPassword(6)
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "has", has, false)

	has, err = svc.HasPassword(7)
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "has", has, true)
}

func Test_UserSvc_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		userId        int
		avatar        *multipart.FileHeader
		req           types.UpdateUserReq
		expectedUser  types.User
		expectedError error
	}{
		{
			name:          "invalid user id",
			userId:        555,
			avatar:        nil,
			req:           types.UpdateUserReq{},
			expectedUser:  types.User{},
			expectedError: services.ErrUserNotFound,
		},
		{
			name:   "normal update user",
			userId: 7,
			avatar: nil,
			req: types.UpdateUserReq{
				Username:  "newuser-updated",
				AvatarUrl: "new avatar url that shouldn't be updated",
			},
			expectedUser: types.User{
				Id:        7,
				Username:  "newuser-updated",
				Email:     "newuser@gmail.com",
				AvatarUrl: "www.newuser.com/pic",
				IsAdmin:   false,
				Coin:      0,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.UpdateUser(tt.userId, tt.avatar, tt.req)
			asserts.EqualError(t, err, tt.expectedError)
			if err != nil {
				return
			}

			user, err := svc.GetUser(tt.userId)
			asserts.EqualError(t, err, nil)
			asserts.Equal(t, "user", user, tt.expectedUser)
		})
	}
}

func Test_UserSvc_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userId        int
		expectedError error
	}{
		{
			name:          "invalid user id",
			userId:        555,
			expectedError: services.ErrUserNotFound,
		},
		{
			name:          "normal delete user",
			userId:        7,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DeleteUser(tt.userId)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UserSvc_HasAccessToken(t *testing.T) {
	passport, err := svc.SignIn("i.deepaung@gmail.com", "password")
	asserts.EqualError(t, err, nil)

	has, err := svc.HasAccessToken(1, passport.Token.AccessToken)
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "has", has, true)

	has, err = svc.HasAccessToken(1, "some random string here")
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "has", has, false)

	err = svc.SignOut(1, passport.Token.Id)
	asserts.EqualError(t, err, nil)
}

func Test_UserSvc_HasRefreshToken(t *testing.T) {
	passport, err := svc.SignIn("i.deepaung@gmail.com", "password")
	asserts.EqualError(t, err, nil)

	has, err := svc.HasRefreshToken(1, passport.Token.RefreshToken)
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "has", has, true)

	has, err = svc.HasRefreshToken(1, "some random string here")
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "has", has, false)

	err = svc.SignOut(1, passport.Token.Id)
	asserts.EqualError(t, err, nil)
}

func Test_UserSvc_UpdateTokens(t *testing.T) {
	passport, err := svc.SignIn("i.deepaung@gmail.com", "password")
	asserts.EqualError(t, err, nil)

	time.Sleep(1 * time.Second)

	tests := []struct {
		name          string
		userId        int
		refreshToken  string
		assertToken   func(t *testing.T, token types.Token)
		expectedError error
	}{
		{
			name:          "invalid id",
			userId:        1000,
			refreshToken:  passport.Token.RefreshToken,
			assertToken:   func(t *testing.T, token types.Token) {},
			expectedError: services.ErrInvalidRefreshToken,
		},
		{
			name:          "invalid refresh token",
			userId:        1,
			refreshToken:  "some random string here",
			assertToken:   func(t *testing.T, token types.Token) {},
			expectedError: services.ErrInvalidRefreshToken,
		},
		{
			name:         "normal update tokens",
			userId:       1,
			refreshToken: passport.Token.RefreshToken,
			assertToken: func(t *testing.T, token types.Token) {
				asserts.Equal(t, "Token.Id", token.Id, passport.Token.Id)
				asserts.NotEqual(
					t,
					"Token.AccessToken",
					token.AccessToken,
					passport.Token.AccessToken,
				)
				asserts.NotEqual(
					t,
					"Token.RefreshToken",
					token.RefreshToken,
					passport.Token.RefreshToken,
				)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.UpdateTokens(tt.userId, tt.refreshToken)
			time.Sleep(500 * time.Microsecond)
			token, err = svc.UpdateTokens(tt.userId, token.RefreshToken)
			tt.assertToken(t, token)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}

	err = svc.SignOut(1, passport.Token.Id)
	asserts.EqualError(t, err, nil)
}

func Test_UserSvc_ToggleFollow(t *testing.T)     {}
func Test_UserSvc_IsFollowing(t *testing.T)      {}
func Test_UserSvc_Follow(t *testing.T)           {}
func Test_UserSvc_UnFollow(t *testing.T)         {}
func Test_UserSvc_generatePassport(t *testing.T) {}
