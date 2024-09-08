package repositories

import (
	"context"
	"errors"
	"net/http"

	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

var (
	ErrNoRowsAffected = func(table string) error {
		return httperror.New("table "+table+" no rows affected", http.StatusInternalServerError)
	}

	ErrNotFound = func(name string) error {
		return httperror.New(name+" not found", http.StatusBadRequest)
	}
)

func HandleExecCtx(
	stmt Statement,
	ctx context.Context,
	db qrm.DB,
	table string,
) error {
	if _, err := RawStatement("PRAGMA foreign_keys=on;").ExecContext(ctx, db); err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNoRowsAffected(table)
	}

	return nil
}

func HandleExecCtxWithErr(
	stmt Statement,
	ctx context.Context,
	db qrm.DB,
	noRowsAffectedErr error,
) error {
	if _, err := RawStatement("PRAGMA foreign_keys=on;").ExecContext(ctx, db); err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return noRowsAffectedErr
	}

	return nil
}

func HandleQueryCtx(
	stmt Statement,
	ctx context.Context,
	db qrm.DB,
	dest interface{},
	name string,
) error {
	if _, err := RawStatement("PRAGMA foreign_keys=on;").ExecContext(ctx, db); err != nil {
		return err
	}

	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return ErrNotFound(name)
		}
		return err
	}

	return nil
}

func HandleQueryCtxWithErr(
	stmt Statement,
	ctx context.Context,
	db qrm.DB,
	dest interface{},
	notFoundErr error,
) error {
	if _, err := RawStatement("PRAGMA foreign_keys=on;").ExecContext(ctx, db); err != nil {
		return err
	}

	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return notFoundErr
		}
		return err
	}

	return nil
}

func HandleHasCtx(
	stmt Statement,
	ctx context.Context,
	db qrm.DB,
	dest interface{},
) (bool, error) {
	if _, err := RawStatement("PRAGMA foreign_keys=on;").ExecContext(ctx, db); err != nil {
		return false, err
	}

	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
