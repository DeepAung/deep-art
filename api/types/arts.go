package types

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/DeepAung/deep-art/.gen/model"
)

/*
- Update art info
- Update art Files&Cover (case create art)
*/
type FullArtDTO struct {
	Name        string `form:"name"        validate:"required"`
	Description string `form:"description"`
	Price       int    `form:"price"`
	TagsID      []int  `form:"tags"`

	Cover *multipart.FileHeader
	Files []*multipart.FileHeader
}

type ArtDTO struct {
	Name        string `form:"name"        validate:"required"`
	Description string `form:"description"`
	Price       int    `form:"price"`
	TagsID      []int  `form:"tags"`
}

type CreateArtReq struct {
	CreatorId   int
	Name        string
	Description string
	Price       int
	TagsID      []int
}

type UpdateArtInfoReq struct {
	ArtId int

	Name        string
	Description string
	Price       int
	TagsID      []int
}

type UpdateArtFilesReq struct {
	ArtId int

	CoverURL string
	FilesURL []string

	FilesName []string
}

type Art struct {
	model.Arts

	// Users     model.Users `alias:"Creator.*"`
	Creator  Creator `alias:"Creator.*"`
	Files    []model.Files
	Tags     []model.Tags
	TagNames string `alias:"Temp.TagNames"`
	TagIDs   string `alias:"Temp.TagIDs"`

	TotalDownloads   int `alias:"Stats.TotalDownloads"`
	WeeklyDownloads  int `alias:"Stats.WeeklyDownloads"`
	MonthlyDownloads int `alias:"Stats.MonthlyDownloads"`
	YearlyDownloads  int `alias:"Stats.YearlyDownloads"`

	TotalStars   int `alias:"Stats.TotalStars"`
	WeeklyStars  int `alias:"Stats.WeeklyStars"`
	MonthlyStars int `alias:"Stats.MonthlyStars"`
	YearlyStars  int `alias:"Stats.YearlyStars"`
}

func (art *Art) FillTags() error {
	if art.TagIDs == "" {
		return nil
	}

	tagIDs := strings.Split(art.TagIDs, ",")
	tagNames := strings.Split(art.TagNames, ",")
	art.Tags = make([]model.Tags, len(tagIDs))

	for i := range len(tagIDs) {
		id, err := strconv.Atoi(tagIDs[i])
		if err != nil {
			return err
		}
		id32 := int32(id)

		art.Tags[i] = model.Tags{
			ID:   &id32,
			Name: tagNames[i],
		}
	}

	return nil
}

type ManyArtsRes struct {
	Arts  ManyArts
	Total int
}

type ManyArts []struct {
	model.Arts

	Creator  Creator `alias:"Creator.*"`
	Tags     []model.Tags
	TagNames string
	TagIDs   string

	TotalDownloads   int `alias:"Stats.TotalDownloads"`
	WeeklyDownloads  int `alias:"Stats.WeeklyDownloads"`
	MonthlyDownloads int `alias:"Stats.MonthlyDownloads"`
	YearlyDownloads  int `alias:"Stats.YearlyDownloads"`

	TotalStars   int `alias:"Stats.TotalStars"`
	WeeklyStars  int `alias:"Stats.WeeklyStars"`
	MonthlyStars int `alias:"Stats.MonthlyStars"`
	YearlyStars  int `alias:"Stats.YearlyStars"`
}

func (arts ManyArts) FillTags() error {
	for i := range arts {
		if arts[i].TagIDs == "" {
			continue
		}

		tagIDs := strings.Split(arts[i].TagIDs, ",")
		tagNames := strings.Split(arts[i].TagNames, ",")
		arts[i].Tags = make([]model.Tags, len(tagIDs))

		for j := range len(tagIDs) {
			id, err := strconv.Atoi(tagIDs[j])
			if err != nil {
				return err
			}
			id32 := int32(id)

			arts[i].Tags[j] = model.Tags{
				ID:   &id32,
				Name: tagNames[j],
			}
		}
	}

	return nil
}

type ManyArtsDTO struct {
	Search     string `query:"search"     json:"search"`
	Filter     string `query:"filter"     json:"filter"`
	Sort       string `query:"sort"       json:"sort"`
	Pagination string `query:"pagination" json:"pagination"`
}

type ManyArtsReq struct {
	Search     string     `query:"search"     json:"search"`
	Filter     Filter     `query:"filter"     json:"filter"`
	Sort       Sort       `query:"sort"       json:"sort"`
	Pagination Pagination `query:"pagination" json:"pagination"`
}

type Filter struct {
	Tags      []string `query:"tags"      json:"tags"`
	MinPrice  int      `query:"minPrice"  json:"minPrice"  validate:"gte=-1"`
	MaxPrice  int      `query:"maxPrice"  json:"maxPrice"  validate:"gte=-1"`
	ImageExts []string `query:"imageExts" json:"imageExts"`
}

type Sort struct {
	By  string `query:"by"  json:"by"`
	Asc bool   `query:"asc" json:"asc"`
}

type Pagination struct {
	Page  int `query:"page"  json:"page"  validate:"gte=1"`
	Limit int `query:"limit" json:"limit" validate:"gte=1"`
}

type By string

const (
	TotalDownloads   By = "totalDownloads"
	WeeklyDownloads  By = "weeklyDownloads"
	MonthlyDownloads By = "monthlyDownloads"
	YearlyDownloads  By = "yearlyDownloads"
	TotalStars       By = "totalStars"
	WeeklyStars      By = "weeklyStars"
	MonthlyStars     By = "monthlyStars"
	YearlyStars      By = "yearlyStars"
	Price            By = "price"
)
