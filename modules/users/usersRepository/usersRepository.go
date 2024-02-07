package usersRepository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/DeepAung/deep-art/modules/users"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	CreateUser(req *users.RegisterReq, isAdmin bool) (*users.UserPassport, error)
	GetUserById(userId int) (*users.User, error)
	GetUserByEmail(email string) (*users.UserWithPassword, error)
	CreateToken(userId int, accessToken, refreshToken string) (int, error)
	DeleteToken(userId, tokenId int) error
	GetUserByOAuth(social users.SocialEnum, socialId string) (bool, *users.User, error)
	HasOAuth(req *users.OAuthReq) bool
	CreateOAuth(req *users.OAuthCreateReq) error
	DeleteOAuth(req *users.OAuthReq) error
	GetTokenInfo(refreshToken string) (*users.TokenInfo, error)
	UpdateToken(token *users.Token) error
	GetUserEmailById(userId int) (string, error)
}

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) CreateUser(
	req *users.RegisterReq,
	isAdmin bool,
) (*users.UserPassport, error) {
	query := `
	INSERT INTO "users" (
		"username",
		"email",
		"password",
    "is_admin"
	)
	VALUES
		($1, $2, $3, $4)
	RETURNING
    "id", "username", "email", "avatar_url", "is_admin";`

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("create user failed: %v", err)
	}

	user := new(users.User)
	err = tx.
		QueryRowx(query, req.Username, req.Email, req.Password, isAdmin).
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

	passport := &users.UserPassport{
		User:  user,
		Token: nil,
	}
	return passport, nil
}

func (r *usersRepository) GetUserById(userId int) (*users.User, error) {
	query := `
  SELECT
    "id",
    "username",
    "email",
    "avatar_url",
    "is_admin"
  FROM "users"
  WHERE "id" = $1
  LIMIT 1`

	user := new(users.User)
	err := r.db.Get(user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
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
    "avatar_url",
    "is_admin"
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

func (r *usersRepository) GetUserByOAuth(
	social users.SocialEnum,
	socialId string,
) (bool, *users.User, error) {
	query := `
  SELECT 
    "u"."id",
    "u"."username",
    "u"."email",
    "u"."avatar_url",
    "u"."is_admin"
  FROM "oauths" AS "o"
  LEFT JOIN "users" AS "u"
  ON "o"."user_id" = "u"."id"
  WHERE "social" = $1 AND "social_id" = $2
  LIMIT 1;`

	user := new(users.User)
	err := r.db.Get(user, query, social, socialId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf("get user by oauth failed: %v", err)
	}

	return true, user, nil
}

func (r *usersRepository) HasOAuth(req *users.OAuthReq) bool {
	query := `
	 SELECT count(*)
	 FROM "oauths"
	 WHERE "user_id" = $1 AND "social" = $2
	 LIMIT 1;`

	var count int
	err := r.db.Get(&count, query, req.UserId, req.Social)
	if err != nil || count == 0 {
		return false
	}

	return true
}

func (r *usersRepository) CreateOAuth(req *users.OAuthCreateReq) error {
	query := `
  INSERT INTO "oauths" (
    "user_id",
    "social",
    "social_id"
  )
  VALUES
    ($1, $2, $3);`

	_, err := r.db.Exec(query, req.UserId, req.Social, req.SocialId)
	if err != nil {
		return fmt.Errorf("create oauth failed: %v", err)
	}

	return nil
}

func (r *usersRepository) DeleteOAuth(req *users.OAuthReq) error {
	query := `
  DELETE FROM "oauths"
  WHERE "user_id" = $1 AND "social" = $2;`

	result, err := r.db.Exec(query, req.UserId, req.Social)
	if err != nil {
		return fmt.Errorf("delete oauth failed: %v", err)
	}

	num, err := result.RowsAffected()
	if err != nil || num == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

func (r *usersRepository) DeleteToken(userId, tokenId int) error {
	query := `
  DELETE FROM "tokens"
  WHERE "id" = $1 AND "user_id" = $2;`

	result, err := r.db.Exec(query, tokenId, userId)
	if err != nil {
		return fmt.Errorf("delete token failed: %v", err)
	}

	num, err := result.RowsAffected()
	if err != nil || num == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

func (r *usersRepository) GetTokenInfo(refreshToken string) (*users.TokenInfo, error) {
	query := `
  SELECT
    "id",
    "user_id"
  FROM "tokens"
  WHERE "refresh_token" = $1
  LIMIT 1;`

	token := new(users.TokenInfo)
	err := r.db.Get(token, query, refreshToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *usersRepository) UpdateToken(token *users.Token) error {
	query := `
  UPDATE "tokens" SET
    "access_token" = $1,
    "refresh_token" = $2
  WHERE "id" = $3;`

	_, err := r.db.Exec(query, token.AccessToken, token.RefreshToken, token.Id)
	if err != nil {
		return fmt.Errorf("update token failed: %v", err)
	}

	return nil
}

func (r *usersRepository) GetUserEmailById(userId int) (string, error) {
	query := `
  SELECT
    "email"
  FROM "users"
  WHERE "id" = $1
  LIMIT 1`

	var email string
	err := r.db.Get(&email, query, userId)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	return email, nil
}
