package repositories

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

var ErrArtsNotFound = httperror.New("art not found", http.StatusBadRequest)

type ArtsRepo struct {
	storer  storer.Storer
	db      *sql.DB
	timeout time.Duration
}

func NewArtsRepo(storer storer.Storer, db *sql.DB, timeout time.Duration) *ArtsRepo {
	return &ArtsRepo{
		storer:  storer,
		db:      db,
		timeout: timeout,
	}
}

func (r *ArtsRepo) FindManyArts(page int) (types.ManyArts, error) {
	infoTable := SELECT(
		Arts.ID,
		COUNT(DISTINCT(DownloadedArts.ID)).AS("TotalDownloads"),
		COUNT(DISTINCT(UsersStarredArts.UserID)).AS("TotalStars"),
	).FROM(
		Arts.
			LEFT_JOIN(DownloadedArts, DownloadedArts.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(UsersStarredArts, UsersStarredArts.ArtID.EQ(Arts.ID)),
	).AsTable("Info")

	creator := Users.AS("Creator")
	cover := Files.AS("Cover")

	stmt := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		cover.AllColumns,
		Files.AllColumns,
		Tags.AllColumns,
		infoTable.AllColumns().As("Info.*"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.AS("Creator").ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(cover, cover.ID.EQ(Arts.CoverID)).
			LEFT_JOIN(Files, Files.ArtID.EQ(Arts.ID).
				AND(Files.ID.NOT_EQ(Arts.CoverID))).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(DownloadedArts, DownloadedArts.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(infoTable, Arts.ID.From(infoTable).EQ(Arts.ID)),
	)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest types.ManyArts
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return types.ManyArts{}, ErrArtsNotFound
		}
		return types.ManyArts{}, err
	}

	return dest, nil
}

func (r *ArtsRepo) FindOneArt(id int) (types.Art, error) {
	infoTable := SELECT(
		Arts.ID,
		COUNT(DISTINCT(DownloadedArts.ID)).AS("TotalDownloads"),
		r.countInterval(DownloadedArts.ID, DownloadedArts.CreatedAt, DAYS(-7)).
			AS("WeeklyDownloads"),
		r.countInterval(DownloadedArts.ID, DownloadedArts.CreatedAt, MONTHS(-1)).
			AS("MonthlyDownloads"),
		r.countInterval(DownloadedArts.ID, DownloadedArts.CreatedAt, YEARS(-1)).
			AS("YearlyDownloads"),

		COUNT(DISTINCT(UsersStarredArts.UserID)).AS("TotalStars"),
		r.countInterval(UsersStarredArts.UserID, UsersStarredArts.CreatedAt, DAYS(-7)).
			AS("WeeklyStars"),
		r.countInterval(UsersStarredArts.UserID, UsersStarredArts.CreatedAt, MONTHS(-1)).
			AS("MonthlyStars"),
		r.countInterval(UsersStarredArts.UserID, UsersStarredArts.CreatedAt, YEARS(-1)).
			AS("YearlyStars"),
	).FROM(
		Arts.
			LEFT_JOIN(DownloadedArts, DownloadedArts.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(UsersStarredArts, UsersStarredArts.ArtID.EQ(Arts.ID)),
	).WHERE(Arts.ID.EQ(Int(int64(id)))).AsTable("Info")

	creator := Users.AS("Creator")
	cover := Files.AS("Cover")

	stmt := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		cover.AllColumns,
		Files.AllColumns,
		Tags.AllColumns,
		infoTable.AllColumns().As("Info.*"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.AS("Creator").ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(cover, cover.ID.EQ(Arts.CoverID)).
			LEFT_JOIN(Files, Files.ArtID.EQ(Arts.ID).
				AND(Files.ID.NOT_EQ(Arts.CoverID))).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(DownloadedArts, DownloadedArts.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(infoTable, Arts.ID.From(infoTable).EQ(Arts.ID)),
	).WHERE(Arts.ID.EQ(Int(int64(id))))

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest types.Art
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return types.Art{}, ErrArtsNotFound
		}
		return types.Art{}, err
	}

	return dest, nil
}

// func (r *ArtsRepo) UploadFiles(
// 	coverFiles *multipart.FileHeader,
// 	artFiles []*multipart.FileHeader,
// 	dir string,
// ) ([]storer.FileRes, error) {
// 	res, err := r.storer.UploadFiles(artFiles, dir)
// 	if err != nil {
// 		return nil, nil
// 	}
//
// 	_ = res
// 	// stmt := Files.
// 	// 	INSERT(Files.ArtID, Files.Filename, Files.Filetype, Files.URL).
// 	// 	VALUES()
// 	// ID        *int32 `sql:"primary_key"`
// 	// ArtID     int32
// 	// Filename  string
// 	// Filetype  string
// 	// URL       string
// 	// CreatedAt *time.Time
// 	// UpdatedAt *time.Time
//
// 	return nil, nil
// }
//
// func (r *ArtsRepo) DeleteFiles(destinations []string) error {
// 	return nil
// }

func (r *ArtsRepo) countInterval(
	id ColumnInteger,
	timestamp TimestampExpression,
	interval Expression,
) Expression {
	return COUNT(DISTINCT(CASE().WHEN(timestamp.GT(DATETIME("now", interval))).THEN(id)))
}
