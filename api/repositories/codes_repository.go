package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
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
	ErrCodeNotFound       = httperror.New("code not found", http.StatusBadRequest)
	ErrCodeNoRowsAffected = httperror.New("code no rows affected", http.StatusInternalServerError)
)

type CodesRepo struct {
	db      *sql.DB
	timeout time.Duration
}

func NewCodesRepo(db *sql.DB, timeout time.Duration) *CodesRepo {
	return &CodesRepo{
		db:      db,
		timeout: timeout,
	}
}

func (r *CodesRepo) FindAllCodes() ([]model.Codes, error) {
	stmt := SELECT(Codes.AllColumns).FROM(Codes)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest []model.Codes
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		return []model.Codes{}, err
	}

	return dest, nil
}

func (r *CodesRepo) FindOneCodeById(id int) (model.Codes, error) {
	stmt := SELECT(Codes.AllColumns).FROM(Codes).WHERE(Codes.ID.EQ(Int(int64(id)))).LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Codes
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return model.Codes{}, ErrCodeNotFound
		}
		return model.Codes{}, err
	}

	log.Println("dest: ", dest)

	return dest, nil
}

func (r *CodesRepo) CreateCode(req types.CodeReq) error {
	stmt := Codes.INSERT(Codes.Name, Codes.Value, Codes.ExpTime).
		VALUES(req.Name, req.Value, req.ExpTime.Time)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrCodeNoRowsAffected
	}

	return nil
}

func (r *CodesRepo) UpdateCode(id int, req types.CodeReq) error {
	stmt := Codes.UPDATE(Codes.Name, Codes.Value, Codes.ExpTime).
		SET(req.Name, req.Value, req.ExpTime.Time).
		WHERE(Codes.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrCodeNoRowsAffected
	}

	return nil
}

func (r *CodesRepo) DeleteCode(id int) error {
	stmt := Codes.DELETE().WHERE(Codes.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrCodeNoRowsAffected
	}

	return nil
}
