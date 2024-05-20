package repositories

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

var (
	ErrUserNotFound       = httperror.New("user not found", http.StatusBadRequest)
	ErrUserNoRowsAffected = httperror.New("user no rows affected", http.StatusInternalServerError)
	ErrCreateUserFailed   = httperror.New("create user failed", http.StatusInternalServerError)
	ErrTokenNotFound      = httperror.New("token not found", http.StatusBadRequest)
)

type UsersRepo struct {
	db      *sql.DB
	timeout time.Duration
}

func NewUsersRepo(db *sql.DB, timeout time.Duration) *UsersRepo {
	return &UsersRepo{
		db:      db,
		timeout: timeout,
	}
}

func (r *UsersRepo) FindOneUserById(id int) (types.User, error) {
	stmt := SELECT(Users.AllColumns).
		FROM(Users).
		WHERE(Users.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Users
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, err
	}

	return types.User{
		Id:        int(*dest.ID),
		Username:  dest.Username,
		Email:     dest.Email,
		AvatarUrl: dest.AvatarURL,
		IsAdmin:   dest.IsAdmin,
		Coin:      int(dest.Coin),
	}, nil
}

func (r *UsersRepo) FindOneUserWithPasswordByEmail(email string) (types.UserWithPassword, error) {
	stmt := SELECT(Users.AllColumns).
		FROM(Users).
		WHERE(Users.Email.EQ(String(email)))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Users
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return types.UserWithPassword{}, ErrUserNotFound
		}
		return types.UserWithPassword{}, err
	}

	return types.UserWithPassword{
		Id:        int(*dest.ID),
		Username:  dest.Username,
		Email:     dest.Email,
		Password:  dest.Password,
		AvatarUrl: dest.AvatarURL,
		IsAdmin:   dest.IsAdmin,
		Coin:      int(dest.Coin),
	}, nil
}

func (r *UsersRepo) CreateUser(req types.SignUpReq) (types.User, error) {
	stmt := Users.INSERT(Users.Username, Users.Email, Users.Password, Users.AvatarURL).
		VALUES(req.Username, req.Email, req.Password, req.AvatarUrl).
		RETURNING(Users.AllColumns)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Users
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		return types.User{}, err
	}

	return types.User{
		Id:        int(*dest.ID),
		Username:  dest.Username,
		Email:     dest.Email,
		AvatarUrl: dest.AvatarURL,
		IsAdmin:   dest.IsAdmin,
		Coin:      int(dest.Coin),
	}, nil
}

func (r *UsersRepo) UpdateUser(id int, req types.UpdateReq) error {
	stmt := Users.UPDATE(Users.Username, Users.Email, Users.AvatarURL).
		SET(req.Username, req.Email, req.AvatarUrl).
		WHERE(Users.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := stmt.ExecContext(ctx, r.db)
	n, rowerr := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 || rowerr != nil {
		return ErrUserNoRowsAffected
	}

	return nil
}

func (r *UsersRepo) DeleteUser(id int) error {
	stmt := Users.DELETE().
		WHERE(Users.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := stmt.ExecContext(ctx, r.db)
	n, rowerr := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 || rowerr != nil {
		return ErrUserNoRowsAffected
	}

	return nil
}

func (r *UsersRepo) HasAccessToken(userId int, accessToken string) (bool, error) {
	stmt := SELECT(Int(1)).
		FROM(Tokens).
		WHERE(
			Tokens.UserID.EQ(Int(int64(userId))).
				AND(Tokens.AccessToken.EQ(String(accessToken))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp struct{ int }
	if err := stmt.QueryContext(ctx, r.db, &tmp); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *UsersRepo) FindOneTokenId(userId int, refreshToken string) (int, error) {
	stmt := SELECT(Tokens.ID).
		FROM(Tokens).
		WHERE(
			Tokens.UserID.EQ(Int(int64(userId))).
				AND(Tokens.RefreshToken.EQ(String(refreshToken))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var res model.Tokens
	if err := stmt.QueryContext(ctx, r.db, &res); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, ErrTokenNotFound
		}
		return 0, err
	}

	return int(*res.ID), nil
}

func (r *UsersRepo) CreateToken(userId int, accessToken, refreshToken string) (id int, err error) {
	stmt := Tokens.INSERT(Tokens.UserID, Tokens.AccessToken, Tokens.RefreshToken).
		VALUES(userId, accessToken, refreshToken).
		RETURNING(Tokens.ID)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var res model.Tokens
	if err = stmt.QueryContext(ctx, r.db, &res); err != nil {
		return 0, err
	}

	return int(*res.ID), nil
}

func (r *UsersRepo) DeleteToken(userId int, tokenId int) error {
	stmt := Tokens.DELETE().WHERE(
		Tokens.UserID.EQ(Int(int64(userId))).
			AND(Tokens.ID.EQ(Int(int64(tokenId)))),
	)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := stmt.ExecContext(ctx, r.db)
	n, rowerr := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 || rowerr != nil {
		return ErrTokenNotFound
	}

	return nil
}
