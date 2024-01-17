package usersRepository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/DeepAung/deep-art/modules/users"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	CreateUser(req *users.RegisterReq) (*users.User, error)
	GetUserByEmail(email string) (*users.UserWithPassword, error)
	CreateToken(userId int, accessToken, refreshToken string) (int, error)
	GetUserIdByOAuth(social users.SocialEnum, socialId string) (bool, int, error)
}

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) CreateUser(req *users.RegisterReq) (*users.User, error) {
	query := `
	INSERT INTO "users" (
		"username",
		"email",
		"password"
	)
	VALUES
		($1, $2, $3)
	RETURNING
    "id", "username", "email", "avatar_url";`

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("create user failed: %v", err)
	}

	user := new(users.User)
	err = tx.
		QueryRowx(query, req.Username, req.Email, req.Password).
		StructScan(user)
	if err != nil {
		tx.Rollback()

		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("create user failed: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("create user failed: %v", err)
	}

	return user, nil
}

func (r *usersRepository) GetUserByEmail(email string) (*users.UserWithPassword, error) {
	query := `
  SELECT
    "id",
    "username",
    "email",
    "password",
    "avatar_url"
  FROM "users"
  WHERE "email" = $1
  LIMIT 1`

	user := new(users.UserWithPassword)
	err := r.db.Get(user, query, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (r *usersRepository) CreateToken(userId int, accessToken, refreshToken string) (int, error) {
	query := `
	INSERT INTO "tokens" (
    "user_id",
		"access_token",
		"refresh_token"
	)
	VALUES
		($1, $2, $3)
	RETURNING
    "id";`

	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("create token failed: %v", err)
	}

	var tokenId int
	err = tx.QueryRow(query, userId, accessToken, refreshToken).Scan(&tokenId)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("create token failed: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("create token failed: %v", err)
	}

	return tokenId, nil
}

func (r *usersRepository) GetUserIdByOAuth(
	social users.SocialEnum,
	socialId string,
) (bool, int, error) {
	query := `
  SELECT "user_id" FROM "oauths"
  WHERE "social" = $1 AND "social_id" = $2
  LIMIT 1;`

	var userId int
	err := r.db.Get(&userId, query, social, socialId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, 0, nil
		}

		return false, 0, fmt.Errorf("get user id by oauth failed: %v", err)
	}

	return true, userId, nil
}
