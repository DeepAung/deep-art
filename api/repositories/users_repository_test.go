package repositories_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/asserts"
	"github.com/golang-migrate/migrate/v4"
)

var (
	testDB    *sql.DB
	migrateDB *migrate.Migrate
	repo      *repositories.UsersRepo
)

func init() {
	testDB, migrateDB = repositories.NewTestDB()
	repositories.ResetDB(migrateDB)
	repo = repositories.NewUsersRepo(testDB, 1*time.Second)
}

func Test_UsersRepo_CreateUser(t *testing.T) {
	normalInput := types.SignUpReq{
		Username:        "newuser",
		Email:           "newuser@gmail.com",
		Password:        "",
		ConfirmPassword: "",
		AvatarUrl:       "",
	}

	uniqueEmailInput := normalInput
	uniqueEmailInput.Email = "i.deepaung@gmail.com"

	uniqueUsernameInput := normalInput
	uniqueUsernameInput.Username = "DeepAung"

	tests := []struct {
		name          string
		input         types.SignUpReq
		expectedUser  types.User
		expectedError error
	}{
		{
			name:          "unique email",
			input:         uniqueEmailInput,
			expectedError: repositories.ErrUniqueEmail,
		},
		{
			name:          "unique username",
			input:         uniqueUsernameInput,
			expectedError: repositories.ErrUniqueUsername,
		},
		{
			name:  "normal user",
			input: normalInput,
			expectedUser: types.User{
				Id:        6,
				Username:  normalInput.Username,
				Email:     normalInput.Email,
				AvatarUrl: normalInput.AvatarUrl,
				IsAdmin:   false,
				Coin:      0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.CreateUser(tt.input)
			asserts.EqualError(t, err, tt.expectedError)
			if err == nil {
				asserts.Equal(t, "user", user, tt.expectedUser)
			}
		})
	}
}

func Test_UsersRepo_FindOneUserById(t *testing.T) {
	tests := []struct {
		name          string
		inputId       int
		expectedUser  types.User
		expectedError error
	}{
		{
			name:    "find by exist id",
			inputId: 1,
			expectedUser: types.User{
				Id:        1,
				Username:  "DeepAung",
				Email:     "i.deepaung@gmail.com",
				AvatarUrl: "",
				IsAdmin:   false,
				Coin:      0,
			},
			expectedError: nil,
		},
		{
			name:          "find by non-exist id",
			inputId:       100,
			expectedUser:  types.User{},
			expectedError: repositories.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindOneUserById(tt.inputId)

			asserts.EqualError(t, err, tt.expectedError)
			asserts.Equal(t, "user", user, tt.expectedUser)
		})
	}
}

func Test_UsersRepo_FindOneUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		inputEmail    string
		expectedUser  types.User
		expectedError error
	}{
		{
			name:       "find by exist email",
			inputEmail: "i.deepaung@gmail.com",
			expectedUser: types.User{
				Id:        1,
				Username:  "DeepAung",
				Email:     "i.deepaung@gmail.com",
				AvatarUrl: "",
				IsAdmin:   false,
				Coin:      0,
			},
			expectedError: nil,
		},
		{
			name:          "find by non-exist email",
			inputEmail:    "non-exist@gmail.com",
			expectedUser:  types.User{},
			expectedError: repositories.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindOneUserByEmail(tt.inputEmail)

			asserts.EqualError(t, err, tt.expectedError)
			asserts.Equal(t, "user", user, tt.expectedUser)
		})
	}
}

func Test_UsersRepo_FindOneUserWithPasswordByEmail(t *testing.T) {
	tests := []struct {
		name          string
		inputEmail    string
		expectedError error
	}{
		{
			name:          "find by exist email",
			inputEmail:    "i.deepaung@gmail.com",
			expectedError: nil,
		},
		{
			name:          "find by non-exist email",
			inputEmail:    "non-exist@gmail.com",
			expectedError: repositories.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindOneUserByEmail(tt.inputEmail)

			asserts.EqualError(t, err, tt.expectedError)
			if (err == nil) && (user == types.User{}) {
				t.Fatalf("expected to have user but got %v", types.User{})
			}
		})
	}
}

func Test_UsersRepo_FindOneCreatorById(t *testing.T) {
	tests := []struct {
		name          string
		inputId       int
		expectedError error
	}{
		{
			name:          "find by exist id",
			inputId:       1,
			expectedError: nil,
		},
		{
			name:          "find by non-exist id",
			inputId:       100,
			expectedError: repositories.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindOneCreatorById(tt.inputId)

			asserts.EqualError(t, err, tt.expectedError)
			if (err == nil) && (user == types.Creator{}) {
				t.Fatalf("expected to have creator but got %v", types.User{})
			}
		})
	}
}

func Test_UsersRepo_HasPassword(t *testing.T) {
	tests := []struct {
		userId        int
		expectedHas   bool
		expectedError error
	}{
		{1, true, nil},
		{2, true, nil},
		{3, true, nil},
		{4, true, nil},
		{5, true, nil},
		{6, false, nil},
		{1000, false, repositories.ErrUserNotFound},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("HasPassword(%d)", tt.userId), func(t *testing.T) {
			has, err := repo.HasPassword(tt.userId)
			asserts.EqualError(t, err, tt.expectedError)
			asserts.Equal(t, "has", has, tt.expectedHas)
		})
	}
}

func Test_UsersRepo_UpdateUser(t *testing.T) {
	normalReq := types.UpdateUserReq{
		Username:  "updated newuser",
		AvatarUrl: "updated avatar",
	}

	uniqueUsernameReq := normalReq
	uniqueUsernameReq.Username = "admin"

	emptyUsernameReq := normalReq
	emptyUsernameReq.Username = ""

	tests := []struct {
		name          string
		inputId       int
		inputReq      types.UpdateUserReq
		expectedError error
	}{
		{
			name:          "update by non-exist id",
			inputId:       1000,
			inputReq:      normalReq,
			expectedError: repositories.ErrNoRowsAffected("users"),
		},
		{
			name:          "update unique username",
			inputId:       6,
			inputReq:      uniqueUsernameReq,
			expectedError: repositories.ErrUniqueUsername,
		},
		{
			name:          "update empty username",
			inputId:       6,
			inputReq:      emptyUsernameReq,
			expectedError: repositories.ErrEmptyUsername,
		},
		{
			name:          "update by exist id",
			inputId:       6,
			inputReq:      normalReq,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateUser(tt.inputId, tt.inputReq)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UsersRepo_UpdateUserPassword(t *testing.T) {
	inputId := 6
	inputPassword := "password02"

	err := repo.UpdateUserPassword(inputId, inputPassword)
	asserts.EqualError(t, err, nil)
}

func Test_UsersRepo_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		inputId       int
		expectedError error
	}{
		{
			name:          "delete by non-exist id",
			inputId:       1000,
			expectedError: repositories.ErrNoRowsAffected("users"),
		},
		{
			name:          "delete by exist id",
			inputId:       6,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteUser(tt.inputId)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UsersRepo_CreateFollow(t *testing.T) {
	tests := []struct {
		name          string
		followerId    int
		followeeId    int
		expectedError error
	}{
		{
			name:          "create non-exist follower_id",
			followerId:    1000,
			followeeId:    1,
			expectedError: repositories.ErrNotExistFollowerFolloweeId,
		},
		{
			name:          "create non-exist followee_id",
			followerId:    1,
			followeeId:    1000,
			expectedError: repositories.ErrNotExistFollowerFolloweeId,
		},
		{
			name:          "unique follower_id & followee_id",
			followerId:    2,
			followeeId:    1,
			expectedError: repositories.ErrUniqueFollow,
		},
		{
			name:          "create normal follow",
			followerId:    5,
			followeeId:    1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateFollow(tt.followerId, tt.followeeId)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UsersRepo_HasFollow(t *testing.T) {
	tests := []struct {
		followerId  int
		followeeId  int
		expectedHas bool
	}{
		{2, 1, true},
		{3, 1, true},
		{4, 1, true},
		{5, 1, true},
		{4, 2, false},
		{2, 3, false},
		{1000, 1000, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("HasFollow(%d, %d)", tt.followerId, tt.followeeId), func(t *testing.T) {
			has, err := repo.HasFollow(tt.followerId, tt.followeeId)
			asserts.EqualError(t, err, nil)
			asserts.Equal(t, "has", has, tt.expectedHas)
		})
	}
}

func Test_UsersRepo_DeleteFollow(t *testing.T) {
	tests := []struct {
		name          string
		followerId    int
		followeeId    int
		expectedError error
	}{
		{
			name:          "delete non-exist follow",
			followerId:    1000,
			followeeId:    1,
			expectedError: repositories.ErrNoRowsAffected("follow"),
		},
		{
			name:          "delete non-exist follow (02)",
			followerId:    1,
			followeeId:    1000,
			expectedError: repositories.ErrNoRowsAffected("follow"),
		},
		{
			name:          "delete normal follow",
			followerId:    5,
			followeeId:    1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteFollow(tt.followerId, tt.followeeId)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UserRepo_CreateToken(t *testing.T) {
	tests := []struct {
		name          string
		userId        int
		accessToken   string
		refreshToken  string
		expectedError error
	}{
		{
			name:          "create by non-exist user_id",
			userId:        1000,
			accessToken:   "example access token",
			refreshToken:  "example refresh token",
			expectedError: repositories.ErrNotExistUserId,
		},
		{
			name:          "create by exist user_id",
			userId:        1,
			accessToken:   "example access token",
			refreshToken:  "example refresh token",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreateToken(tt.userId, tt.accessToken, tt.refreshToken)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UsersRepo_HasAccessToken(t *testing.T) {
	tests := []struct {
		userId      int
		accessToken string
		expectedHas bool
	}{
		{1, "example access token", true},
		{2, "example access token", false},
		{1000, "example access token", false},
		{1, "random access token", false},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("HasAccessToken(%d, %s)", tt.userId, tt.accessToken)
		t.Run(name, func(t *testing.T) {
			has, err := repo.HasAccessToken(tt.userId, tt.accessToken)
			asserts.EqualError(t, err, nil)
			asserts.Equal(t, "has", has, tt.expectedHas)
		})
	}
}

func Test_UsersRepo_HasRefreshToken(t *testing.T) {
	tests := []struct {
		userId       int
		refreshToken string
		expectedHas  bool
	}{
		{1, "example refresh token", true},
		{2, "example refresh token", false},
		{1000, "example refresh token", false},
		{1, "random refresh token", false},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("HasRefreshToken(%d, %s)", tt.userId, tt.refreshToken)
		t.Run(name, func(t *testing.T) {
			has, err := repo.HasRefreshToken(tt.userId, tt.refreshToken)
			asserts.EqualError(t, err, nil)
			asserts.Equal(t, "has", has, tt.expectedHas)
		})
	}
}

func Test_UsersRepo_FindOneTokenId(t *testing.T) {
	tests := []struct {
		name          string
		userId        int
		refreshToken  string
		expectedId    int
		expectedError error
	}{
		{
			name:         "find by non-exist user_id",
			userId:       1000,
			refreshToken: "example refresh token",
			// expectedId:    0,
			expectedError: repositories.ErrTokenNotFound,
		},
		{
			name:         "find by non-exist refreshToken",
			userId:       1,
			refreshToken: "random refresh token",
			// expectedId:    0,
			expectedError: repositories.ErrTokenNotFound,
		},
		{
			name:          "find by exist user_id & refreshToken",
			userId:        1,
			refreshToken:  "example refresh token",
			expectedId:    1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.FindOneTokenId(tt.userId, tt.refreshToken)
			asserts.EqualError(t, err, tt.expectedError)
			if err == nil {
				asserts.Equal(t, "token id", id, tt.expectedId)
			}
		})
	}
}

func Test_UsersRepo_UpdateTokens(t *testing.T) {
	tests := []struct {
		name          string
		tokenId       int
		accessToken   string
		refreshToken  string
		expectedError error
	}{
		{
			name:          "update by non-exist id",
			tokenId:       1000,
			accessToken:   "updated access token",
			refreshToken:  "updated refresh token",
			expectedError: repositories.ErrNoRowsAffected("tokens"),
		},
		{
			name:          "update by exist id",
			tokenId:       1,
			accessToken:   "updated access token",
			refreshToken:  "updated refresh token",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateTokens(tt.tokenId, tt.accessToken, tt.refreshToken)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}

func Test_UsersRepo_DeleteToken(t *testing.T) {
	tests := []struct {
		name          string
		userId        int
		tokenId       int
		expectedError error
	}{
		{
			name:          "delete by non-exist user id",
			userId:        1000,
			tokenId:       1,
			expectedError: repositories.ErrNoRowsAffected("tokens"),
		},
		{
			name:          "delete by non-exist token id",
			userId:        1,
			tokenId:       1000,
			expectedError: repositories.ErrNoRowsAffected("tokens"),
		},
		{
			name:          "delete by exist token",
			userId:        1,
			tokenId:       1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteToken(tt.userId, tt.tokenId)
			asserts.EqualError(t, err, tt.expectedError)
		})
	}
}
