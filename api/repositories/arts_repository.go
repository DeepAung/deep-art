package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	ErrInvalidSortingType = httperror.New("invalid sorting type", http.StatusBadRequest)
	ErrInvalidPrice       = httperror.New(
		"invalid price, please try again",
		http.StatusBadRequest,
	)
	ErrArtsTagsNoRowsAffected = ErrNoRowsAffected("arts_tags")
	ErrFilesNoRowsAffected    = ErrNoRowsAffected("files")
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

func (r *ArtsRepo) BeginTx() (context.Context, context.CancelFunc, *sql.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	tx, err := r.db.BeginTx(ctx, nil)

	return ctx, cancel, tx, err
}

func (r *ArtsRepo) CreateArt(req types.CreateArtReq) (artId int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	artId, err = r.CreateArtWithDB(ctx, r.db, req)
	return
}

func (r *ArtsRepo) CreateArtWithDB(
	ctx context.Context,
	db qrm.DB,
	req types.CreateArtReq,
) (artId int, err error) {
	var art struct {
		ID int `alias:"Arts.ID"`
	}
	// insert arts
	stmt1 := Arts.
		INSERT(Arts.Name, Arts.Description, Arts.CreatorID, Arts.Price, Arts.CoverURL).
		VALUES(req.Name, req.Description, req.CreatorId, req.Price, "unset").
		RETURNING(Arts.ID)
	if err = HandleQueryCtx(stmt1, ctx, db, &art, "art"); err != nil {
		return
	}

	fmt.Println("artId: ", art.ID)

	// insert tags
	if req.TagsID != nil && len(req.TagsID) > 0 {
		fmt.Println("insert tags")
		stmt2 := ArtsTags.INSERT(ArtsTags.ArtID, ArtsTags.TagID)
		for _, tagID := range req.TagsID {
			stmt2 = stmt2.VALUES(art.ID, tagID)
		}
		if err = HandleExecCtx(stmt2, ctx, db, "arts_tags"); err != nil {
			return
		}
	}

	artId = art.ID
	return
}

func (r *ArtsRepo) UpdateArtInfo(req types.UpdateArtInfoReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.UpdateArtInfoWithDB(ctx, r.db, req)
}

func (r *ArtsRepo) UpdateArtInfoWithDB(
	ctx context.Context,
	db qrm.DB,
	req types.UpdateArtInfoReq,
) error {
	stmt1 := Arts.
		UPDATE(Arts.Name, Arts.Description, Arts.Price).
		SET(req.Name, req.Description, req.Price).
		WHERE(Arts.ID.EQ(Int(int64(req.ArtId))))
	if err := HandleExecCtx(stmt1, ctx, db, "arts"); err != nil {
		return err
	}

	stmt2 := ArtsTags.DELETE().WHERE(ArtsTags.ArtID.EQ(Int(int64(req.ArtId))))
	if err := HandleExecCtxWithErr(stmt2, ctx, db, ErrArtsTagsNoRowsAffected); err != nil &&
		!errors.Is(err, ErrArtsTagsNoRowsAffected) {
		return err
	}

	if req.TagsID != nil && len(req.TagsID) > 0 {
		stmt3 := ArtsTags.INSERT(ArtsTags.ArtID, ArtsTags.TagID)
		for _, tagID := range req.TagsID {
			stmt3 = stmt3.VALUES(req.ArtId, tagID)
		}
		if err := HandleExecCtx(stmt3, ctx, db, "arts_tags"); err != nil {
			return err
		}
	}

	return nil
}

func (r *ArtsRepo) DeleteArtWithDB(ctx context.Context, db qrm.DB, artId int) error {
	stmt := Arts.DELETE().WHERE(Arts.ID.EQ(Int(int64(artId))))

	return HandleExecCtx(stmt, ctx, db, "arts")
}

func (r *ArtsRepo) FindManyArts(req types.ManyArtsReq) (types.ManyArtsRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	statsTable := r.statsTable().AsTable("Stats")
	stats := r.statsColumn(statsTable)
	creator := Users.AS("Creator")

	var cond BoolExpression = Int(1).EQ(Int(1))
	cond = r.withFilterCond(cond, req.Filter)
	cond = r.withSearchCond(cond, req.Search)
	cond = r.withCreatorIdCond(cond, req.CreatorId)

	// stmt
	stmt := SELECT(
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
	stmt, err := r.withSortStmt(stmt, req.Sort, stats)
	if err != nil {
		return types.ManyArtsRes{}, err
	}
	stmt = r.withPaginationStmt(stmt, req.Pagination)

	// stmt2
	stmt2 := SELECT(
		COUNT(DISTINCT(Arts.ID)).AS("Count"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond)

	// query stmt
	var dest types.ManyArts
	if err := HandleQueryCtx(stmt, ctx, r.db, &dest, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}
	err = dest.FillTags()
	if err != nil {
		return types.ManyArtsRes{}, err
	}

	// query stmt2
	var dest2 struct{ Count int }
	if err := HandleQueryCtx(stmt2, ctx, r.db, &dest2, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}

	return types.ManyArtsRes{
		Arts:  dest,
		Total: dest2.Count,
	}, nil
}

func (r *ArtsRepo) FindManyStarredArts(
	userId int,
	req types.ManyArtsReq,
) (types.ManyArtsRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	statsTable := r.statsTable().AsTable("Stats")
	stats := r.statsColumn(statsTable)
	creator := Users.AS("Creator")

	var cond BoolExpression = Int(1).EQ(Int(1))
	cond = r.withFilterCond(cond, req.Filter)
	cond = r.withSearchCond(cond, req.Search)

	// starred arts stmt
	stmt := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		COUNT(DISTINCT(Follow.UserIDFollower)).AS("Creator.Followers"),
		Raw("group_concat(DISTINCT tags.name)").AS("TagNames"),
		Raw("group_concat(DISTINCT tags.id)").AS("TagIDs"),
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			INNER_JOIN(
				UsersStarredArts,
				UsersStarredArts.ArtID.EQ(Arts.ID).
					AND(UsersStarredArts.UserID.EQ(Int(int64(userId)))),
			).
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond).GROUP_BY(Arts.ID)
	stmt, err := r.withSortStmt(stmt, req.Sort, stats)
	if err != nil {
		return types.ManyArtsRes{}, err
	}
	stmt = r.withPaginationStmt(stmt, req.Pagination)

	// starred arts stmt2
	stmt2 := SELECT(
		COUNT(DISTINCT(Arts.ID)).AS("Count"),
	).FROM(
		Arts.
			INNER_JOIN(
				UsersStarredArts,
				UsersStarredArts.ArtID.EQ(Arts.ID).
					AND(UsersStarredArts.UserID.EQ(Int(int64(userId)))),
			).
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond)

	// query stmt
	var dest types.ManyArts
	if err := HandleQueryCtx(stmt, ctx, r.db, &dest, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}
	err = dest.FillTags()
	if err != nil {
		return types.ManyArtsRes{}, err
	}

	// query stmt2
	var dest2 struct{ Count int }
	if err := HandleQueryCtx(stmt2, ctx, r.db, &dest2, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}

	return types.ManyArtsRes{
		Arts:  dest,
		Total: dest2.Count,
	}, nil
}

func (r *ArtsRepo) FindManyBoughtArts(
	userId int,
	req types.ManyArtsReq,
) (types.ManyArtsRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	statsTable := r.statsTable().AsTable("Stats")
	stats := r.statsColumn(statsTable)
	creator := Users.AS("Creator")

	var cond BoolExpression = Int(1).EQ(Int(1))
	cond = r.withFilterCond(cond, req.Filter)
	cond = r.withSearchCond(cond, req.Search)

	// bought arts stmt
	stmt := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		COUNT(DISTINCT(Follow.UserIDFollower)).AS("Creator.Followers"),
		Raw("group_concat(DISTINCT tags.name)").AS("TagNames"),
		Raw("group_concat(DISTINCT tags.id)").AS("TagIDs"),
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			INNER_JOIN(
				UsersBoughtArts,
				UsersBoughtArts.ArtID.EQ(Arts.ID).
					AND(UsersBoughtArts.UserID.EQ(Int(int64(userId)))),
			).
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond).GROUP_BY(Arts.ID)
	stmt, err := r.withSortStmt(stmt, req.Sort, stats)
	if err != nil {
		return types.ManyArtsRes{}, err
	}
	stmt = r.withPaginationStmt(stmt, req.Pagination)

	// bought arts stmt2
	stmt2 := SELECT(
		COUNT(DISTINCT(Arts.ID)).AS("Count"),
	).FROM(
		Arts.
			INNER_JOIN(
				UsersBoughtArts,
				UsersBoughtArts.ArtID.EQ(Arts.ID).
					AND(UsersBoughtArts.UserID.EQ(Int(int64(userId)))),
			).
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond)

	// query stmt
	var dest types.ManyArts
	if err := HandleQueryCtx(stmt, ctx, r.db, &dest, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}
	err = dest.FillTags()
	if err != nil {
		return types.ManyArtsRes{}, err
	}

	// query stmt2
	var dest2 struct{ Count int }
	if err := HandleQueryCtx(stmt2, ctx, r.db, &dest2, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}

	return types.ManyArtsRes{
		Arts:  dest,
		Total: dest2.Count,
	}, nil
}

func (r *ArtsRepo) FindManyCreatedArts(
	userId int,
	req types.ManyArtsReq,
) (types.ManyArtsRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	statsTable := r.statsTable().AsTable("Stats")
	stats := r.statsColumn(statsTable)
	creator := Users.AS("Creator")

	var cond BoolExpression = Int(1).EQ(Int(1))
	cond = r.withFilterCond(cond, req.Filter)
	cond = r.withSearchCond(cond, req.Search)

	// created arts stmt
	stmt := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		COUNT(DISTINCT(Follow.UserIDFollower)).AS("Creator.Followers"),
		Raw("group_concat(DISTINCT tags.name)").AS("TagNames"),
		Raw("group_concat(DISTINCT tags.id)").AS("TagIDs"),
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			INNER_JOIN(
				creator,
				creator.ID.EQ(Arts.CreatorID).
					AND(creator.ID.EQ(Int(int64(userId)))),
			).
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond).GROUP_BY(Arts.ID)
	stmt, err := r.withSortStmt(stmt, req.Sort, stats)
	if err != nil {
		return types.ManyArtsRes{}, err
	}
	stmt = r.withPaginationStmt(stmt, req.Pagination)

	// created arts stmt2
	stmt2 := SELECT(
		COUNT(DISTINCT(Arts.ID)).AS("Count"),
	).FROM(
		Arts.
			INNER_JOIN(
				creator,
				creator.ID.EQ(Arts.CreatorID).
					AND(creator.ID.EQ(Int(int64(userId)))),
			).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(cond)

	// query stmt
	var dest types.ManyArts
	if err := HandleQueryCtx(stmt, ctx, r.db, &dest, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}
	err = dest.FillTags()
	if err != nil {
		return types.ManyArtsRes{}, err
	}

	// query stmt2
	var dest2 struct{ Count int }
	if err := HandleQueryCtx(stmt2, ctx, r.db, &dest2, "art"); err != nil {
		return types.ManyArtsRes{}, err
	}

	return types.ManyArtsRes{
		Arts:  dest,
		Total: dest2.Count,
	}, nil
}

func (r *ArtsRepo) FindOneArt(id int) (types.Art, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	statsTable := r.statsTable().
		WHERE(Arts.ID.EQ(Int(int64(id)))).
		AsTable("Stats")

	creator := Users.AS("Creator")

	stmt1 := SELECT(
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		COUNT(DISTINCT(Follow.UserIDFollower)).AS("Creator.Followers"),
		Raw("group_concat(DISTINCT tags.name)").AS("Temp.TagNames"),
		Raw("group_concat(DISTINCT tags.id)").AS("Temp.TagIDs"),
		statsTable.AllColumns().As("Stats.*"),
	).FROM(
		Arts.
			LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID)),
	).WHERE(Arts.ID.EQ(Int(int64(id))))

	var dest types.Art
	if err := HandleQueryCtx(stmt1, ctx, r.db, &dest, "art"); err != nil {
		return types.Art{}, err
	}

	stmt2 := SELECT(Arts.ID, Files.AllColumns).
		FROM(Arts.LEFT_JOIN(Files, Files.ArtID.EQ(Arts.ID))).
		WHERE(Arts.ID.EQ(Int(int64(id))))

	var filesDest types.Art
	if err := HandleQueryCtx(stmt2, ctx, r.db, &filesDest, "art"); err != nil {
		return types.Art{}, err
	}

	dest.Files = filesDest.Files
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
	return HandleHasCtx(stmt, ctx, r.db, &tmp)
}

func (r *ArtsRepo) FindUserCoin(userId int) (int, error) {
	stmt := SELECT(Users.Coin).FROM(Users).WHERE(Users.ID.EQ(Int(int64(userId))))
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user model.Users
	if err := HandleQueryCtx(stmt, ctx, r.db, &user, "user"); err != nil {
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
	if err := HandleExecCtx(stmt1, ctx, tx, "users_bought_arts"); err != nil {
		return err
	}

	stmt2 := SELECT(Arts.Price).FROM(Arts).WHERE(Arts.ID.EQ(Int(int64(artId))))
	var art model.Arts
	if err := HandleQueryCtx(stmt2, ctx, tx, &art, "art"); err != nil {
		return err
	}

	if art.Price != int32(price) {
		return ErrInvalidPrice
	}

	stmt3 := Users.UPDATE(Users.Coin).
		SET(Users.Coin.SUB(Int(int64(art.Price)))).
		WHERE(Users.ID.EQ(Int(int64(userId))))
	if err := HandleExecCtx(stmt3, ctx, tx, "users"); err != nil {
		return err
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
	return HandleHasCtx(stmt, ctx, r.db, &tmp)
}

func (r *ArtsRepo) CreateUsersStarredArts(userId, artId int) error {
	stmt := UsersStarredArts.
		INSERT(UsersStarredArts.UserID, UsersStarredArts.ArtID).
		VALUES(userId, artId)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "users_starred_arts")
}

func (r *ArtsRepo) DeleteUsersStarredArts(userId, artId int) error {
	stmt := UsersStarredArts.DELETE().WHERE(
		UsersStarredArts.UserID.EQ(Int(int64(userId))).
			AND(UsersStarredArts.ArtID.EQ(Int(int64(artId)))),
	)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return HandleExecCtx(stmt, ctx, r.db, "users_starred_arts")
}

func (r *ArtsRepo) FindCreatorID(artId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	stmt := SELECT(Arts.CreatorID).FROM(Arts).WHERE(Arts.ID.EQ(Int(int64(artId))))
	var art model.Arts
	if err := HandleQueryCtx(stmt, ctx, r.db, &art, "art"); err != nil {
		return 0, err
	}

	return int(art.CreatorID), nil
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

func (r *ArtsRepo) findManyArtsColumns(statsTable SelectTable) []Projection {
	creator := Users.AS("Creator")

	return []Projection{
		Arts.AllColumns,
		creator.AllColumns.Except(creator.Password),
		COUNT(DISTINCT(Follow.UserIDFollower)).AS("Creator.Followers"),
		Raw("group_concat(DISTINCT tags.name)").AS("TagNames"),
		Raw("group_concat(DISTINCT tags.id)").AS("TagIDs"),
		statsTable.AllColumns().As("Stats.*"),
	}
}

func (r *ArtsRepo) findManyArtsTable(statsTable SelectTable) ReadableTable {
	creator := Users.AS("Creator")

	return Arts.
		LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
		LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
		LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
		LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
		LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID))
}

func (r *ArtsRepo) countManyArtsTable(statsTable SelectTable) ReadableTable {
	creator := Users.AS("Creator")

	return Arts.
		LEFT_JOIN(creator, creator.ID.EQ(Arts.CreatorID)).
		LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
		LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
		LEFT_JOIN(statsTable, Arts.ID.From(statsTable).EQ(Arts.ID))
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
			LEFT_JOIN(ArtsTags, ArtsTags.ArtID.EQ(Arts.ID)).
			LEFT_JOIN(Tags, Tags.ID.EQ(ArtsTags.TagID)).
			LEFT_JOIN(Follow, Follow.UserIDFollowee.EQ(creator.ID)).
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

func (r *ArtsRepo) withCreatorIdCond(cond BoolExpression, creatorId int) BoolExpression {
	if creatorId == 0 {
		return cond
	}

	return cond.AND(Arts.CreatorID.EQ(Int(int64(creatorId))))
}

// func (r *ArtsRepo) withStarredArtsCond(cond BoolExpression) BoolExpression {
// 	return cond.AND(Arts.ID.BETWEEN)
// }

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
