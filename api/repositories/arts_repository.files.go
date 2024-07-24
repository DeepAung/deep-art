package repositories

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func (r *ArtsRepo) FindOneFile(id int) (model.Files, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	stmt := SELECT(Files.AllColumns).FROM(Files).WHERE(Files.ID.EQ(Int(int64(id))))

	var res model.Files
	err := HandleQueryCtx(stmt, ctx, r.db, &res, "file")
	return res, err
}

func (r *ArtsRepo) FindOneCoverURL(artId int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	stmt := SELECT(Arts.CoverURL).FROM(Arts).WHERE(Arts.ID.EQ(Int(int64(artId))))

	var res model.Arts
	err := HandleQueryCtx(stmt, ctx, r.db, &res, "art")
	return res.CoverURL, err
}

func (r *ArtsRepo) UpdateArtCoverAndFilesWithDB(
	ctx context.Context,
	db qrm.DB,
	req types.UpdateArtFilesReq,
) error {
	stmt1 := Arts.UPDATE(Arts.CoverURL).
		SET(req.CoverURL).
		WHERE(Arts.ID.EQ(Int(int64(req.ArtId))))
	if err := HandleExecCtx(stmt1, ctx, db, "arts"); err != nil {
		return err
	}

	stmt2 := Files.DELETE().WHERE(Files.ArtID.EQ(Int(int64(req.ArtId))))
	if err := HandleExecCtxWithErr(stmt2, ctx, db, ErrFilesNoRowsAffected); err != nil &&
		!errors.Is(err, ErrFilesNoRowsAffected) {
		return err
	}

	stmt3 := Files.INSERT(Files.ArtID, Files.Filename, Files.Filetype, Files.URL)
	for i := range req.FilesURL {
		fileUrl := req.FilesURL[i]
		filename := req.FilesName[i]
		ext := filepath.Ext(filename)

		stmt3 = stmt3.VALUES(req.ArtId, filename, ext, fileUrl)
	}
	if err := HandleExecCtx(stmt3, ctx, db, "files"); err != nil {
		return err
	}

	return nil
}

func (r *ArtsRepo) InsertArtFilesWithDB(
	ctx context.Context,
	db qrm.DB,
	artId int,
	filesURL, filesName []string,
) error {
	stmt := Files.INSERT(Files.ArtID, Files.Filename, Files.Filetype, Files.URL)
	for i := range filesURL {
		fileUrl := filesURL[i]
		filename := filesName[i]
		ext := filepath.Ext(filename)

		stmt = stmt.VALUES(artId, filename, ext, fileUrl)
	}

	return HandleExecCtx(stmt, ctx, db, "files")
}

func (r *ArtsRepo) DeleteArtFilesWithDB(ctx context.Context, db qrm.DB, fileId int) error {
	stmt := Files.DELETE().WHERE(Files.ID.EQ(Int(int64(fileId))))
	return HandleExecCtxWithErr(stmt, ctx, db, ErrFilesNoRowsAffected)
}

func (r *ArtsRepo) UpdateArtCoverWithDB(
	ctx context.Context,
	db qrm.DB,
	artId int,
	coverURL string,
) error {
	stmt := Arts.UPDATE(Arts.CoverURL).SET(coverURL).WHERE(Arts.ID.EQ(Int(int64(artId))))
	return HandleExecCtx(stmt, ctx, db, "arts")
}
