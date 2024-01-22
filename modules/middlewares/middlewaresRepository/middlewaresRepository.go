package middlewaresRepository

import "github.com/jmoiron/sqlx"

type IMiddlewaresRepository interface {
	FindAccessToken(userId int, accessToken string) bool
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func NewMiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId int, accessToken string) bool {
	query := `
  SELECT 1 FROM "tokens"
  WHERE "user_id" = $1 AND "access_token" = $2
  LIMIT 1;`

	var tmp bool
	return r.db.Get(&tmp, query, userId, accessToken) == nil
}
