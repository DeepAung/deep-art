package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	. "github.com/go-jet/jet/v2/sqlite"
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

func (r *TagsRepo) CreateTag(name string) error {
	stmt := Tags.INSERT(Tags.Name).VALUES(name)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "tags")
}

func (r *TagsRepo) UpdateTag(id int, name string) error {
	stmt := Tags.UPDATE(Tags.Name).SET(name).WHERE(Tags.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "tags")
}

func (r *TagsRepo) DeleteTag(id int) error {
	stmt := Tags.DELETE().WHERE(Tags.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "tags")
}
