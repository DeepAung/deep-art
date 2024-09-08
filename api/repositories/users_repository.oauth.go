package repositories

import (
	"context"
	"net/http"

	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

var (
	ErrAlreadyConnect = httperror.New("you already connect to this OAuth", http.StatusBadRequest)
	ErrNotConnectYet  = httperror.New("you did not connect to this OAuth", http.StatusBadRequest)
)

func (r *UsersRepo) FindAllOAuthProvider(userId int) ([]types.OAuthProvider, error) {
	stmt := SELECT(Oauths.Provider).
		FROM(Oauths).
		WHERE(Oauths.UserID.EQ(Int(int64(userId))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var oauths []types.OAuthProvider
	err := HandleQueryCtx(stmt, ctx, r.db, &oauths, "oauth")

	return oauths, err
}

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

	err := HandleExecCtx(stmt, ctx, db, "oauths")
	if err != nil && err.Error() == "UNIQUE constraint failed: oauths.provider, oauths.provider_user_id" {
		return ErrAlreadyConnect
	}
	return err
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

	err := HandleExecCtx(stmt, ctx, db, "oauths")
	if err != nil && err.Error() == "500: table oauths no rows affected" {
		return ErrNotConnectYet
	}
	return err
}
