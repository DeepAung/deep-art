package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/api/types"
	. "github.com/go-jet/jet/v2/sqlite"
)

var (
	ErrUserNotFound  = ErrNotFound("user")
	ErrTokenNotFound = ErrNotFound("token")
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
		WHERE(Users.ID.EQ(Int(int64(id)))).
		LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Users
	if err := HandleQueryCtx(stmt, ctx, r.db, &dest, "user"); err != nil {
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

func (r *UsersRepo) FindOneCreatorById(id int) (types.Creator, error) {
	stmt := SELECT(
		Users.AllColumns,
		COUNT(Follow.UserIDFollower).AS("Followers"),
	).FROM(
		Users.
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(Users.ID)),
	).WHERE(Users.ID.EQ(Int(int64(id)))).GROUP_BY(Users.ID)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest struct {
		model.Users
		Followers int
	}
	if err := HandleQueryCtx(stmt, ctx, r.db, &dest, "user"); err != nil {
		return types.Creator{}, err
	}

	return types.Creator{
		Id:        int(*dest.ID),
		Username:  dest.Username,
		Email:     dest.Email,
		AvatarURL: dest.AvatarURL,
		Followers: dest.Followers,
	}, nil
}

func (r *UsersRepo) FindOneUserWithPasswordByEmail(email string) (types.UserWithPassword, error) {
	stmt := SELECT(Users.AllColumns).
		FROM(Users).
		WHERE(Users.Email.EQ(String(email))).
		LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Users
	if err := HandleQueryCtxWithErr(stmt, ctx, r.db, &dest, ErrUserNotFound); err != nil {
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
	if err := HandleQueryCtx(stmt, ctx, r.db, &dest, "user"); err != nil {
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

	return HandleExecCtx(stmt, ctx, r.db, "users")
}

func (r *UsersRepo) DeleteUser(id int) error {
	stmt := Users.DELETE().
		WHERE(Users.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "users")
}

func (r *UsersRepo) HasFollow(followerId, followeeId int) (bool, error) {
	stmt := SELECT(Int(1)).
		FROM(Follow).
		WHERE(
			Follow.UserIDFollower.EQ(Int(int64(followerId))).
				AND(Follow.UserIDFollowee.EQ(Int(int64(followeeId)))),
		)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp struct{ int }
	return HandleHasCtx(stmt, ctx, r.db, &tmp)
}

func (r *UsersRepo) CreateFollow(followerId, followeeId int) error {
	stmt := Follow.
		INSERT(Follow.UserIDFollower, Follow.UserIDFollowee).
		VALUES(followerId, followeeId)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "follow")
}

func (r *UsersRepo) DeleteFollow(followerId, followeeId int) error {
	stmt := Follow.DELETE().WHERE(
		Follow.UserIDFollower.EQ(Int(int64(followerId))).
			AND(Follow.UserIDFollowee.EQ(Int(int64(followeeId)))),
	)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "follow")
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
	return HandleHasCtx(stmt, ctx, r.db, &tmp)
}

func (r *UsersRepo) HasRefreshToken(userId int, refreshToken string) (bool, error) {
	stmt := SELECT(Int(1)).
		FROM(Tokens).
		WHERE(
			Tokens.UserID.EQ(Int(int64(userId))).
				AND(Tokens.RefreshToken.EQ(String(refreshToken))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp struct{ int }
	return HandleHasCtx(stmt, ctx, r.db, &tmp)
}

func (r *UsersRepo) FindOneTokenId(userId int, refreshToken string) (int, error) {
	stmt := SELECT(Tokens.ID).
		FROM(Tokens).
		WHERE(
			Tokens.UserID.EQ(Int(int64(userId))).
				AND(Tokens.RefreshToken.EQ(String(refreshToken)))).
		LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var res model.Tokens
	if err := HandleQueryCtxWithErr(stmt, ctx, r.db, &res, ErrTokenNotFound); err != nil {
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
	if err := HandleQueryCtx(stmt, ctx, r.db, &res, "token"); err != nil {
		return 0, err
	}

	return int(*res.ID), nil
}

func (r *UsersRepo) UpdateTokens(tokenId int, newAccessToken, newRefreshToken string) error {
	stmt := Tokens.UPDATE(Tokens.AccessToken, Tokens.RefreshToken).
		SET(newAccessToken, newRefreshToken).
		WHERE(Tokens.ID.EQ(Int(int64(tokenId))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "tokens")
}

func (r *UsersRepo) DeleteToken(userId int, tokenId int) error {
	stmt := Tokens.DELETE().WHERE(
		Tokens.UserID.EQ(Int(int64(userId))).
			AND(Tokens.ID.EQ(Int(int64(tokenId)))),
	)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "tokens")
}
