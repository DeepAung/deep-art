package repositories

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	. "github.com/DeepAung/deep-art/.gen/table"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

var (
	ErrArtsNotFound       = httperror.New("art not found", http.StatusBadRequest)
	ErrInvalidSortingType = httperror.New("invalid sorting type", http.StatusBadRequest)

	ErrStarNoRowsAffected = httperror.New(
		"users_starred_arts no rows affected",
		http.StatusInternalServerError,
	)
	ErrBoughtNoRowsAffected = httperror.New(
		"users_bought_arts no rows affected",
		http.StatusInternalServerError,
	)
	ErrInvalidPrice = httperror.New(
		"invalid price, please try again",
		http.StatusBadRequest,
	)
)

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

func (r *ArtsRepo) FindManyArts(req types.ManyArtsReq) (types.ManyArtsRes, error) {
	statsTable := r.statsTable().AsTable("Stats")
	stats := r.statsColumn(statsTable)

	var cond BoolExpression = Arts.ID.EQ(Arts.ID)
	cond = r.withFilterCond(cond, req.Filter)
	cond = r.withSearchCond(cond, req.Search)

	stmt := r.findManyArtsStmt(cond, statsTable)
	stmt, err := r.withSortStmt(stmt, req.Sort, stats)
	if err != nil {
		return types.ManyArtsRes{}, err
	}
	stmt = r.withPaginationStmt(stmt, req.Pagination)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var dest types.ManyArts
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return types.ManyArtsRes{}, ErrArtsNotFound
		}
		return types.ManyArtsRes{}, err
	}

	err = dest.FillTags()
	if err != nil {
		return types.ManyArtsRes{}, err
	}

	stmt2 := r.findCountManyArtsStmt(cond, statsTable)
	var dest2 struct{ Count int }
	if err := stmt2.QueryContext(ctx, r.db, &dest2); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return types.ManyArtsRes{}, ErrArtsNotFound
		}
		return types.ManyArtsRes{}, err
	}

	return types.ManyArtsRes{
		Arts:  dest,
		Total: dest2.Count,
	}, nil
}

func (r *ArtsRepo) FindOneArt(id int) (types.Art, error) {
	statsTable := r.statsTable().
		WHERE(Arts.ID.EQ(Int(int64(id)))).
		AsTable("Stats")

	creator := Users.AS("Creator")

	stmt := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		COUNT(DISTINCT(Follow.UserIDFollower)).AS("Creator.Followers"),
		Files.AllColumns,
		Raw("group_concat(DISTINCT tags.name)").AS("Temp.TagNames"),
		Raw("group_concat(DISTINCT tags.id)").AS("Temp.TagIDs"),
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
			LEFT_JOIN(Files, Files.ArtID.EQ(Arts.ID)).
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

	err := dest.FillTags()
	if err != nil {
		return types.Art{}, err
	}

	return dest, nil
}

func (r *ArtsRepo) HasUsersBoughtArts(userId, artId int) (bool, error) {
	stmt := SELECT(Int(1)).
		FROM(UsersBoughtArts).
		WHERE(
			UsersBoughtArts.UserID.EQ(Int(int64(userId))).
				AND(UsersBoughtArts.ArtID.EQ(Int(int64(artId)))),
		)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp struct{ int }
	if err := stmt.QueryContext(ctx, r.db, &tmp); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *ArtsRepo) FindUserCoin(userId int) (int, error) {
	stmt := SELECT(Users.Coin).FROM(Users).WHERE(Users.ID.EQ(Int(int64(userId))))
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user model.Users
	if err := stmt.QueryContext(ctx, r.db, &user); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, err
	}

	return int(user.Coin), nil
}

func (r *ArtsRepo) BuyArt(userId, artId, price int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt1 := UsersBoughtArts.
		INSERT(UsersBoughtArts.UserID, UsersBoughtArts.ArtID).
		VALUES(userId, artId)
	result, err := stmt1.ExecContext(ctx, tx)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrBoughtNoRowsAffected
	}

	stmt2 := SELECT(Arts.Price).FROM(Arts).WHERE(Arts.ID.EQ(Int(int64(artId))))
	var art model.Arts
	if err := stmt2.QueryContext(ctx, tx, &art); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return ErrArtsNotFound
		}
		return err
	}

	if art.Price != int32(price) {
		return ErrInvalidPrice
	}

	stmt3 := Users.UPDATE(Users.Coin).
		SET(Users.Coin.SUB(Int(int64(art.Price)))).
		WHERE(Users.ID.EQ(Int(int64(userId))))
	result, err = stmt3.ExecContext(ctx, tx)
	if err != nil {
		return err
	}
	n, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUserNoRowsAffected
	}

	return tx.Commit()
}

func (r *ArtsRepo) HasUsersStarredArts(userId, artId int) (bool, error) {
	stmt := SELECT(Int(1)).
		FROM(UsersStarredArts).
		WHERE(
			UsersStarredArts.UserID.EQ(Int(int64(userId))).
				AND(UsersStarredArts.ArtID.EQ(Int(int64(artId)))),
		)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp struct{ int }
	if err := stmt.QueryContext(ctx, r.db, &tmp); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *ArtsRepo) CreateUsersStarredArts(userId, artId int) error {
	stmt := UsersStarredArts.
		INSERT(UsersStarredArts.UserID, UsersStarredArts.ArtID).
		VALUES(userId, artId)

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
		return ErrStarNoRowsAffected
	}

	return nil
}

func (r *ArtsRepo) DeleteUsersStarredArts(userId, artId int) error {
	stmt := UsersStarredArts.DELETE().WHERE(
		UsersStarredArts.UserID.EQ(Int(int64(userId))).
			AND(UsersStarredArts.ArtID.EQ(Int(int64(artId)))),
	)

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
		return ErrStarNoRowsAffected
	}

	return nil
}

// ---------------------------------------------- //

type statsColumn struct {
	totalDownloads   Column
	weeklyDownloads  Column
	monthlyDownloads Column
	yearlyDownloads  Column
	totalStars       Column
	weeklyStars      Column
	monthlyStars     Column
	yearlyStars      Column
}

func (r *ArtsRepo) statsTable() SelectStatement {
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

	return SELECT(
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
	)
}

func (r *ArtsRepo) statsColumn(statsTable SelectTable) statsColumn {
	return statsColumn{
		totalDownloads:   IntegerColumn("TotalDownloads").From(statsTable),
		weeklyDownloads:  IntegerColumn("WeeklyDownloads").From(statsTable),
		monthlyDownloads: IntegerColumn("MonthlyDownloads").From(statsTable),
		yearlyDownloads:  IntegerColumn("YearlyDownloads").From(statsTable),
		totalStars:       IntegerColumn("TotalStars").From(statsTable),
		weeklyStars:      IntegerColumn("WeeklyStars").From(statsTable),
		monthlyStars:     IntegerColumn("MonthlyStars").From(statsTable),
		yearlyStars:      IntegerColumn("YearlyStars").From(statsTable),
	}
}

func (r *ArtsRepo) findManyArtsStmt(
	cond BoolExpression,
	statsTable SelectTable,
) SelectStatement {
	creator := Users.AS("Creator")

	return SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		COUNT(DISTINCT(Follow.UserIDFollower)).AS("Creator.Followers"),
		Raw("group_concat(DISTINCT tags.name)").AS("TagNames"),
		Raw("group_concat(DISTINCT tags.id)").AS("TagIDs"),
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond).GROUP_BY(Arts.ID)
}

func (r *ArtsRepo) findCountManyArtsStmt(
	cond BoolExpression,
	statsTable SelectTable,
) SelectStatement {
	creator := Users.AS("Creator")

	return SELECT(
		COUNT(DISTINCT(Arts.ID)).AS("Count"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
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
	stats statsColumn,
) (SelectStatement, error) {
	if sort.By == "" {
		return stmt, nil
	}

	var orderBy Expression
	by := types.By(sort.By)
	switch by {
	case types.TotalDownloads:
		orderBy = stats.totalDownloads
	case types.WeeklyDownloads:
		orderBy = stats.weeklyDownloads
	case types.MonthlyDownloads:
		orderBy = stats.monthlyDownloads
	case types.YearlyDownloads:
		orderBy = stats.yearlyDownloads
	case types.TotalStars:
		orderBy = stats.totalStars
	case types.WeeklyStars:
		orderBy = stats.weeklyStars
	case types.MonthlyStars:
		orderBy = stats.monthlyStars
	case types.YearlyStars:
		orderBy = stats.yearlyStars
	case types.Price:
		orderBy = Arts.Price
	default:
		return nil, ErrInvalidSortingType
	}

	if sort.Asc {
		stmt = stmt.ORDER_BY(orderBy)
	} else {
		stmt = stmt.ORDER_BY(orderBy.DESC())
	}

	return stmt, nil
}

func (r *ArtsRepo) withPaginationStmt(
	stmt SelectStatement,
	pagination types.Pagination,
) SelectStatement {
	page := pagination.Page
	limit := pagination.Limit
	return stmt.LIMIT(int64(limit)).OFFSET(int64(limit*page - limit))
}

func (r *ArtsRepo) countInterval(
	id ColumnInteger,
	timestamp TimestampExpression,
	interval Expression,
) Expression {
	return COUNT(DISTINCT(CASE().WHEN(timestamp.GT(DATETIME("now", interval))).THEN(id)))
}
