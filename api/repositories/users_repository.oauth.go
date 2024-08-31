package repositories

import (
	"context"

	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func (r *UsersRepo) CreateOAuth(userId int, provider string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.CreateOAuthWithDB(ctx, r.db, userId, provider)
}

func (r *UsersRepo) DeleteOAuth(userId int, provider string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.DeleteOAuthWithDB(ctx, r.db, userId, provider)
}

func (r *UsersRepo) CreateOAuthWithDB(
	ctx context.Context,
	db qrm.DB,
	userId int,
	provider string,
) error {
	stmt := Oauths.
		INSERT(Oauths.UserID, Oauths.Provider).
		VALUES(userId, provider)

	err := HandleExecCtx(stmt, ctx, db, "oauths")
	if err == nil {
		return nil
	}

	switch err.Error() {
	case "jet: UNIQUE constraint failed: users.email":
		return ErrUniqueEmail
	case "jet: UNIQUE constraint failed: users.username":
		return ErrUniqueUsername
	default:
		return err
	}
}

// TODO: err handling not found
func (r *UsersRepo) DeleteOAuthWithDB(
	ctx context.Context,
	db qrm.DB,
	userId int,
	provider string,
) error {
	stmt := Oauths.DELETE().
		WHERE(Oauths.UserID.EQ(Int(int64(userId))).AND(Oauths.Provider.EQ(String(provider))))

	return HandleExecCtx(stmt, ctx, db, "oauths")
}
