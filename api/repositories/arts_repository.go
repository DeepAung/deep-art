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
	stats := r.statsExpression()
	statsTable := r.statsTable(stats)

	var cond BoolExpression = Arts.ID.EQ(Arts.ID)
	cond = r.withFilterCond(cond, req.Filter)
	cond = r.withSearchCond(cond, req.Search)

	stmt := r.findManyArtsStmt(cond, statsTable)
	stmt = r.withSortStmt(stmt, req.Sort, stats)
	// TODO:
	// stmt = r.withPaginationStmt(stmt, req.Page)

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
	statsTable := SELECT(
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
	).WHERE(Arts.ID.EQ(Int(int64(id)))).AsTable("Stats")

	creator := Users.AS("Creator")
	cover := Files.AS("Cover")

	stmt := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		cover.AllColumns,
		Files.AllColumns,
		Tags.AllColumns,
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.AS("Creator").ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(cover, cover.ID.EQ(Arts.CoverID)).
			LEFT_JOIN(Files, Files.ArtID.EQ(Arts.ID).
				AND(Files.ID.NOT_EQ(Arts.CoverID))).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
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

// ---------------------------------------------- //

type statsExpression struct {
	totalDownloads   Expression
	weeklyDownloads  Expression
	monthlyDownloads Expression
	yearlyDownloads  Expression
	totalStars       Expression
	weeklyStars      Expression
	monthlyStars     Expression
	yearlyStars      Expression
}

func (r *ArtsRepo) statsExpression() statsExpression {
	stats := statsExpression{}

	id := DownloadedArts.ID
	time := DownloadedArts.CreatedAt

	stats.totalDownloads = COUNT(DISTINCT(id))
	stats.weeklyDownloads = r.countInterval(id, time, DAYS(-7))
	stats.monthlyDownloads = r.countInterval(id, time, MONTHS(-1))
	stats.yearlyDownloads = r.countInterval(id, time, YEARS(-1))

	id = UsersStarredArts.UserID
	time = UsersStarredArts.CreatedAt

	stats.totalStars = COUNT(DISTINCT(id))
	stats.weeklyStars = r.countInterval(id, time, DAYS(-7))
	stats.monthlyStars = r.countInterval(id, time, MONTHS(-1))
	stats.yearlyStars = r.countInterval(id, time, YEARS(-1))

	return stats
}

func (r *ArtsRepo) statsTable(stats statsExpression) SelectTable {
	return SELECT(
		Arts.ID,

		stats.totalDownloads.AS("TotalDownloads"),
		stats.weeklyDownloads.AS("WeeklyDownloads"),
		stats.monthlyDownloads.AS("MonthlyDownloads"),
		stats.yearlyDownloads.AS("YearlyDownloads"),

		stats.totalStars.AS("TotalStars"),
		stats.weeklyStars.AS("WeeklyStars"),
		stats.monthlyStars.AS("MonthlyStars"),
		stats.yearlyStars.AS("YearlyStars"),
	).FROM(
		Arts.
			LEFT_JOIN(DownloadedArts, DownloadedArts.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(UsersStarredArts, UsersStarredArts.ArtID.EQ(Arts.ID)),
	).AsTable("Stats")
}

func (r *ArtsRepo) findManyArtsStmt(
	cond BoolExpression,
	statsTable SelectTable,
) SelectStatement {
	creator := Users.AS("Creator")
	cover := Files.AS("Cover")

	return SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		cover.AllColumns,
		Tags.AllColumns,
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(cover, cover.ID.EQ(Arts.CoverID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond)
}

func (r *ArtsRepo) withFilterCond(cond BoolExpression, filter types.Filter) BoolExpression {
	if filter.Tags != nil && len(filter.Tags) >= 0 {
		strToExp := func(s string) Expression { return String(s) }
		tagsExp := utils.Map(filter.Tags, strToExp)

		artsIDs := SELECT(Arts.ID).
			FROM(
				Arts.
					LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
					LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID).
						AND(Tags.Name.IN(tagsExp...))),
			).GROUP_BY(Arts.ID).
			HAVING(
				COUNT(ArtsTags.TagID).GT_EQ(COUNT(Tags.ID)).
					AND(COUNT(Tags.ID).EQ(Int(int64(len(filter.Tags))))),
			)

		cond = cond.AND(Arts.ID.IN(artsIDs))
	}

	if filter.MinPrice != -1 {
		cond = cond.AND(Arts.Price.GT_EQ(Int(int64(filter.MinPrice))))
	}

	if filter.MaxPrice != -1 {
		cond = cond.AND(Arts.Price.LT_EQ(Int(int64(filter.MaxPrice))))
	}

	// TODO:
	// if filter.ImageExts != nil && len(filter.ImageExts) > 0 {
	// }

	return cond
}

func (r *ArtsRepo) withSearchCond(cond BoolExpression, search string) BoolExpression {
	if search == "" {
		return cond
	}

	creator := Users.AS("Creator")

	return cond.AND(
		Arts.Name.LIKE(String("%" + search + "%")).
			OR(Arts.Description.LIKE(String("%" + search + "%"))).
			OR(creator.Username.LIKE(String("%" + search + "%"))),
	)
}

func (r *ArtsRepo) withSortStmt(
	stmt SelectStatement,
	sort types.Sort,
	stats statsExpression,
) SelectStatement {
	if sort.By == "" {
		return stmt
	}

	var orderBy Expression
	by := types.By(sort.By)
	switch by {
	case types.WeeklyDownloads:
		orderBy = stats.weeklyDownloads
	case types.MonthlyDownloads:
		orderBy = stats.monthlyDownloads
	case types.YearlyDownloads:
		orderBy = stats.yearlyDownloads
	case types.WeeklyStars:
		orderBy = stats.weeklyStars
	case types.MonthlyStars:
		orderBy = stats.monthlyStars
	case types.YearlyStars:
		orderBy = stats.yearlyStars
	case types.Price:
		orderBy = Arts.Price
	default:
		// TODO:
	}

	if sort.Asc {
		stmt = stmt.ORDER_BY(orderBy)
	} else {
		stmt = stmt.ORDER_BY(orderBy.DESC())
	}

	return stmt
}

func (r *ArtsRepo) withPaginationStmt(stmt SelectStatement, page int) SelectStatement {
	const items = 5
	return stmt.LIMIT(items).OFFSET(int64(items*page - items))
}

func (r *ArtsRepo) countInterval(
	id ColumnInteger,
	timestamp TimestampExpression,
	interval Expression,
) Expression {
	return COUNT(DISTINCT(CASE().WHEN(timestamp.GT(DATETIME("now", interval))).THEN(id)))
}
