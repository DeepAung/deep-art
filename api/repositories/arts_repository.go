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
	"github.com/DeepAung/deep-art/pkg/utils"
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

func (r *ArtsRepo) FindManyArts(req types.ManyArtsReq) (types.ManyArts, error) {
	id := DownloadedArts.ID
	time := DownloadedArts.CreatedAt

	totalDownloads := COUNT(DISTINCT(id))
	weeklyDownloads := r.countInterval(id, time, DAYS(-7))
	monthlyDownloads := r.countInterval(id, time, MONTHS(-1))
	yearlyDownloads := r.countInterval(id, time, YEARS(-1))

	id = UsersStarredArts.UserID
	time = UsersStarredArts.CreatedAt

	totalStars := COUNT(DISTINCT(id))
	weeklyStars := r.countInterval(id, time, DAYS(-7))
	monthlyStars := r.countInterval(id, time, MONTHS(-1))
	yearlyStars := r.countInterval(id, time, YEARS(-1))

	// subquery
	infoTable := SELECT(
		Arts.ID,

		totalDownloads.AS("TotalDownloads"),
		weeklyDownloads.AS("WeeklyDownloads"),
		monthlyDownloads.AS("MonthlyDownloads"),
		yearlyDownloads.AS("YearlyDownloads"),

		totalStars.AS("TotalStars"),
		weeklyStars.AS("WeeklyStars"),
		monthlyStars.AS("MonthlyStars"),
		yearlyStars.AS("YearlyStars"),
	).FROM(
		Arts.
			LEFT_JOIN(DownloadedArts, DownloadedArts.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(UsersStarredArts, UsersStarredArts.ArtID.EQ(Arts.ID)),
	).AsTable("Info")

	creator := Users.AS("Creator")
	cover := Files.AS("Cover")

	// query list
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
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(infoTable, Arts.ID.From(infoTable).EQ(Arts.ID)),
	)

	// search
	stmt = stmt.WHERE(
		Arts.Name.LIKE(String("%" + req.Search + "%")).
			OR(Arts.Description.LIKE(String("%" + req.Search + "%"))).
			OR(creator.Username.LIKE(String("%" + req.Search + "%"))),
	)

	// filter
	toStringExp := func(s string) Expression { return String(s) }
	tagsExp := utils.Map(req.Filter.Tags, toStringExp)

	if len(req.Filter.Tags) > 0 {
		stmt = stmt.
			WHERE(Tags.Name.IN(tagsExp...)).
			HAVING(
				COUNT(ArtsTags.TagID).GT_EQ(COUNT(Tags.ID)).
					AND(COUNT(Tags.ID).EQ(Int(int64(len(req.Filter.Tags))))),
			)
	}

	if req.Filter.MinPrice != -1 {
		stmt = stmt.WHERE(Arts.Price.GT_EQ(Int(int64(req.Filter.MinPrice))))
	}

	if req.Filter.MaxPrice != -1 {
		stmt = stmt.WHERE(Arts.Price.GT_EQ(Int(int64(req.Filter.MaxPrice))))
	}

	if len(req.Filter.ImageExts) > 0 {
		// TODO:
	}

	// sort
	if req.Sort.By != "" {
		by := types.By(req.Sort.By)
		switch by {
		case types.WeeklyDownloads:
			stmt = stmt.ORDER_BY(weeklyDownloads)
		case types.MonthlyDownloads:
			stmt = stmt.ORDER_BY(monthlyDownloads)
		case types.YearlyDownloads:
			stmt = stmt.ORDER_BY(yearlyDownloads)
		case types.WeeklyStars:
			stmt = stmt.ORDER_BY(weeklyStars)
		case types.MonthlyStars:
			stmt = stmt.ORDER_BY(monthlyStars)
		case types.YearlyStars:
			stmt = stmt.ORDER_BY(yearlyStars)
		case types.Price:
			stmt = stmt.ORDER_BY(Arts.Price)
		default:
			// TODO:
		}
	}

	// pagination
	stmt = stmt.LIMIT(50).OFFSET(int64(50*req.Page - 50))

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
