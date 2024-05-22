package repositories

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

var (
	ErrTagNotFound       = httperror.New("tag not found", http.StatusBadRequest)
	ErrTagNoRowsAffected = httperror.New("tag no rows affected", http.StatusInternalServerError)
)

type TagsRepo struct {
	db      *sql.DB
	timeout time.Duration
}

func NewTagsRepo(db *sql.DB, timeout time.Duration) *TagsRepo {
	return &TagsRepo{
		db:      db,
		timeout: timeout,
	}
}

func (r *TagsRepo) FindAllTags() ([]model.Tags, error) {
	stmt := SELECT(Tags.AllColumns).FROM(Tags)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest []model.Tags
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		return []model.Tags{}, err
	}

	return dest, nil
}

func (r *TagsRepo) FindOneTagById(id int) (model.Tags, error) {
	stmt := SELECT(Tags.AllColumns).FROM(Tags).WHERE(Tags.ID.EQ(Int(int64(id)))).LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Tags
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return model.Tags{}, ErrTagNotFound
		}
		return model.Tags{}, err
	}

	return dest, nil
}

func (r *TagsRepo) CreateTag(name string) error {
	stmt := Tags.INSERT(Tags.Name).VALUES(name)

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
		return ErrTagNoRowsAffected
	}

	return nil
}

func (r *TagsRepo) UpdateTag(id int, name string) error {
	stmt := Tags.UPDATE(Tags.Name).SET(name).WHERE(Tags.ID.EQ(Int(int64(id))))

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
		return ErrTagNoRowsAffected
	}

	return nil
}

func (r *TagsRepo) DeleteTag(id int) error {
	stmt := Tags.DELETE().WHERE(Tags.ID.EQ(Int(int64(id))))

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
		return ErrTagNoRowsAffected
	}

	return nil
}
