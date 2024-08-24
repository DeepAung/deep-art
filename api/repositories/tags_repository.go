package repositories

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/pkg/httperror"
	. "github.com/go-jet/jet/v2/sqlite"
)

var ErrUniqueTagName = httperror.New("Tag name should be unique", http.StatusBadRequest)

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
	err := HandleQueryCtx(stmt, ctx, r.db, &dest, "tag")
	return dest, err
}

func (r *TagsRepo) FindOneTagById(id int) (model.Tags, error) {
	stmt := SELECT(Tags.AllColumns).FROM(Tags).WHERE(Tags.ID.EQ(Int(int64(id)))).LIMIT(1)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Tags
	err := HandleQueryCtx(stmt, ctx, r.db, &dest, "tag")
	return dest, err
}

func (r *TagsRepo) CreateTag(name string) (model.Tags, error) {
	stmt := Tags.INSERT(Tags.Name).VALUES(name).RETURNING(Tags.AllColumns)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest model.Tags
	err := HandleQueryCtx(stmt, ctx, r.db, &dest, "tag")
	if err != nil && err.Error() == "jet: UNIQUE constraint failed: tags.name" {
		return model.Tags{}, ErrUniqueTagName
	}
	return dest, err
}

func (r *TagsRepo) UpdateTag(id int, name string) error {
	stmt := Tags.UPDATE(Tags.Name).SET(name).WHERE(Tags.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	err := HandleExecCtx(stmt, ctx, r.db, "tags")
	if err != nil && err.Error() == "UNIQUE constraint failed: tags.name" {
		return ErrUniqueTagName
	}
	return err
}

func (r *TagsRepo) DeleteTag(id int) error {
	stmt := Tags.DELETE().WHERE(Tags.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "tags")
}
