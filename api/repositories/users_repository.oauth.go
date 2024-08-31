package repositories

import (
	"context"

	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func (r *UsersRepo) HasOAuth(providerUserId, provider string) (bool, error) {
	stmt := SELECT(Int(1)).
		FROM(Oauths).
		WHERE(Oauths.ProviderUserID.EQ(String(providerUserId)).AND(Oauths.Provider.EQ(String(provider))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp struct{ int }
	return HandleHasCtx(stmt, ctx, r.db, &tmp)
}

func (r *UsersRepo) CreateOAuth(userId int, provider, providerUserId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.CreateOAuthWithDB(ctx, r.db, userId, provider, providerUserId)
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
	providerUserId string,
) error {
	stmt := Oauths.
		INSERT(Oauths.UserID, Oauths.Provider, Oauths.ProviderUserID).
		VALUES(userId, provider, providerUserId)

	return HandleExecCtx(stmt, ctx, db, "oauths")
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
