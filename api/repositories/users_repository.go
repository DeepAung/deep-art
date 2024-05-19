package repositories

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
)

var (
	ErrUserNotFound     = httperror.New("user not found", http.StatusBadGateway)
	ErrCreateUserFailed = httperror.New("create user failed", http.StatusInternalServerError)
	ErrTokenNotFound    = httperror.New("token not found", http.StatusBadGateway)
)

type UsersRepo struct{}

func NewUsersRepo() *UsersRepo {
	return &UsersRepo{}
}

func (r *UsersRepo) FindOneUserById(id int) (types.User, error) {
	return types.User{}, ErrUserNotFound
}
func (r *UsersRepo) FindOneUserWithPasswordByEmail(email string) (types.UserWithPassword, error) {
	return types.UserWithPassword{}, ErrUserNotFound
}

func (r *UsersRepo) CreateUser(req types.SignUpReq) (types.User, error) {
	return types.User{}, ErrCreateUserFailed
}
func (r *UsersRepo) UpdateUser(id int, req types.UpdateReq) error {
	return ErrUserNotFound
}
func (r *UsersRepo) DeleteUser(id int) error { return ErrUserNotFound }

func (r *UsersRepo) FindOneTokenId(userId int, refreshToken string) (int, error) {
	return 0, ErrTokenNotFound
}

func (r *UsersRepo) CreateToken(userId int, accessToken, refreshToken string) (id int, err error) {
	return 0, ErrUserNotFound
}

func (r *UsersRepo) DeleteToken(userId int, tokenId int) error {
	return ErrUserNotFound
}
