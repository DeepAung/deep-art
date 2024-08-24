package repositories

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	. "github.com/go-jet/jet/v2/sqlite"
)

var ErrUniqueCodeName = httperror.New("Code name should be unique", http.StatusBadRequest)

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
	err := HandleQueryCtx(stmt, ctx, r.db, &dest, "code")
	return dest, err
}

func (r *CodesRepo) FindOneCodeById(id int) (model.Codes, error) {
	stmt := SELECT(Codes.AllColumns).FROM(Codes).WHERE(Codes.ID.EQ(Int(int64(id)))).LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Codes
	err := HandleQueryCtx(stmt, ctx, r.db, &dest, "code")
	return dest, err
}

func (r *CodesRepo) FindOneCodeByName(name string) (model.Codes, error) {
	stmt := SELECT(Codes.AllColumns).FROM(Codes).WHERE(Codes.Name.EQ(String(name))).LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Codes
	err := HandleQueryCtx(stmt, ctx, r.db, &dest, "code")
	return dest, err
}

func (r *CodesRepo) HasUsedCode(userId, codeId int) (bool, error) {
	stmt := SELECT(Int(1)).
		FROM(UsersUsedCodes).
		WHERE(
			UsersUsedCodes.UserID.EQ(Int(int64(userId))).
				AND(UsersUsedCodes.CodeID.EQ(Int(int64(codeId)))),
		).
		LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp struct{ int }
	return HandleHasCtx(stmt, ctx, r.db, &tmp)
}

func (r *CodesRepo) UseCode(userId, codeId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt := UsersUsedCodes.
		INSERT(UsersUsedCodes.UserID, UsersUsedCodes.CodeID).
		VALUES(userId, codeId)
	if err := HandleExecCtx(stmt, ctx, tx, "users_used_codes"); err != nil {
		return err
	}

	stmt2 := SELECT(Codes.Value).FROM(Codes).WHERE(Codes.ID.EQ(Int(int64(codeId))))
	var dest2 model.Codes
	if err := HandleQueryCtx(stmt2, ctx, tx, &dest2, "code"); err != nil {
		return err
	}

	stmt3 := Users.
		UPDATE(Users.Coin).
		SET(Users.Coin.ADD(Int(int64(dest2.Value)))).
		WHERE(Users.ID.EQ(Int(int64(userId))))
	if err := HandleExecCtx(stmt3, ctx, tx, "users"); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CodesRepo) CreateCode(req types.CodeReq) (model.Codes, error) {
	stmt := Codes.INSERT(Codes.Name, Codes.Value, Codes.ExpTime).
		VALUES(req.Name, req.Value, req.ExpTime.Time).
		RETURNING(Codes.AllColumns)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Codes
	err := HandleQueryCtx(stmt, ctx, r.db, &dest, "code")
	if err != nil && err.Error() == "jet: UNIQUE constraint failed: codes.name" {
		return model.Codes{}, ErrUniqueCodeName
	}
	return dest, err
}

func (r *CodesRepo) UpdateCode(id int, req types.CodeReq) error {
	stmt := Codes.UPDATE(Codes.Name, Codes.Value, Codes.ExpTime).
		SET(req.Name, req.Value, req.ExpTime.Time).
		WHERE(Codes.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	err := HandleExecCtx(stmt, ctx, r.db, "codes")
	if err != nil && err.Error() == "UNIQUE constraint failed: codes.name" {
		return ErrUniqueCodeName
	}
	return err
}

func (r *CodesRepo) DeleteCode(id int) error {
	stmt := Codes.DELETE().WHERE(Codes.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "codes")
}
