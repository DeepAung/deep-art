package tagsRepository

import (
	"fmt"

	"github.com/DeepAung/deep-art/modules/tags"
	"github.com/jmoiron/sqlx"
)

type ITagsRepository interface {
	GetTags() (*[]tags.Tag, error)
	CreateTag(req *tags.TagReq) error
	UpdateTag(req *tags.TagReq, id int) error
	DeleteTag(id int) error
}

type tagsRepository struct {
	db *sqlx.DB
}

func NewTagsRepository(db *sqlx.DB) ITagsRepository {
	return &tagsRepository{
		db: db,
	}
}
func (r *tagsRepository) GetTags() (*[]tags.Tag, error) {
	query := `SELECT * FROM "tags";`
	tags := new([]tags.Tag)

	err := r.db.Select(tags, query)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagsRepository) CreateTag(req *tags.TagReq) error {
	query := `
  INSERT INTO "tags"
    ("name")
  VALUES
    ($1);`

	_, err := r.db.Exec(query, req.Name)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"tags_name_key\" (SQLSTATE 23505)" {
			return fmt.Errorf("tag name already exists")
		}
		return err
	}

	return nil
}

func (r *tagsRepository) UpdateTag(req *tags.TagReq, id int) error {
	query := `
  UPDATE "tags" SET
    "name" = $1
  WHERE "id" = $2;`

	result, err := r.db.Exec(query, req.Name, id)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"tags_name_key\" (SQLSTATE 23505)" {
			return fmt.Errorf("tag name already exists")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag not found")
	}

	return err
}

func (r *tagsRepository) DeleteTag(id int) error {
	query := `
  DELETE FROM "tags"
  WHERE "id" = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}
