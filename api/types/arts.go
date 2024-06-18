package types

import (
	"strconv"
	"strings"

	"github.com/DeepAung/deep-art/.gen/model"
)

type Art struct {
	model.Arts

	// Users     model.Users `alias:"Creator.*"`
	Creator Creator `alias:"Creator.*"`
	Files   []model.Files
	Tags    []model.Tags

	TotalDownloads   int `alias:"Stats.TotalDownloads"`
	WeeklyDownloads  int `alias:"Stats.WeeklyDownloads"`
	MonthlyDownloads int `alias:"Stats.MonthlyDownloads"`
	YearlyDownloads  int `alias:"Stats.YearlyDownloads"`

	TotalStars   int `alias:"Stats.TotalStars"`
	WeeklyStars  int `alias:"Stats.WeeklyStars"`
	MonthlyStars int `alias:"Stats.MonthlyStars"`
	YearlyStars  int `alias:"Stats.YearlyStars"`
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

type ManyArtsReq struct {
	Search     string     `json:"search"     form:"search"`
	Filter     Filter     `json:"filter"`
	Sort       Sort       `json:"sort"`
	Pagination Pagination `json:"pagination"`
}

type Filter struct {
	Tags      []string `json:"tags"      form:"filter.tags"`
	MinPrice  int      `json:"minPrice"  form:"filter.minPrice"  validate:"gte=-1"`
	MaxPrice  int      `json:"maxPrice"  form:"filter.maxPrice"  validate:"gte=-1"`
	ImageExts []string `json:"imageExts" form:"filter.imageExts"`
}

type Sort struct {
	By  string `json:"by"  form:"sort.by"`
	Asc bool   `json:"asc" form:"sort.asc"`
}

type Pagination struct {
	Page  int `json:"page"  form:"pagination.page"  validate:"gte=1"`
	Limit int `json:"limit" form:"pagination.limit" validate:"gte=1"`
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
