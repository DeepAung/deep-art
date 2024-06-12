package types

import "github.com/DeepAung/deep-art/.gen/model"

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

	Creator model.Users `alias:"Creator.*"`
	Cover   model.Files `alias:"Cover.*"`
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
	WeeklyDownloads  By = "weeklyDownloads"
	MonthlyDownloads By = "monthlyDownloads"
	YearlyDownloads  By = "yearlyDownloads"
	WeeklyStars      By = "weeklyStars"
	MonthlyStars     By = "monthlyStars"
	YearlyStars      By = "yearlyStars"
	Price            By = "price"
)
