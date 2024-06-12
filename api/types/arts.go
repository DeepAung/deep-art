package types

import (
	"strconv"
	"strings"

	"github.com/DeepAung/deep-art/.gen/model"
)

type Art struct {
	model.Arts

	Creator model.Users `alias:"Creator.*"`
	Cover   model.Files `alias:"Cover.*"`
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

type ManyArts []struct {
	model.Arts

	Creator  model.Users `alias:"Creator.*"`
	Cover    model.Files `alias:"Cover.*"`
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
	Search string `form:"search" json:"search"`
	Filter Filter `              json:"filter"`
	Sort   Sort   `              json:"sort"`
	Page   int    `form:"page"   json:"page"`
}

type Filter struct {
	Tags      []string `form:"filter.tags"      json:"tags"`
	MinPrice  int      `form:"filter.minPrice"  json:"minPrice"`
	MaxPrice  int      `form:"filter.maxPrice"  json:"maxPrice"`
	ImageExts []string `form:"filter.imageExts" json:"imageExts"`
}

type Sort struct {
	By  string `form:"sort.by"  json:"by"`
	Asc bool   `form:"sort.asc" json:"asc"`
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
