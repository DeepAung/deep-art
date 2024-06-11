package types

import "github.com/DeepAung/deep-art/.gen/model"

type Art struct {
	model.Arts

	Creator model.Users `alias:"Creator.*"`
	Cover   model.Files `alias:"Cover.*"`
	Files   []model.Files
	Tags    []model.Tags

	TotalDownloads   int `alias:"Info.TotalDownloads"`
	WeeklyDownloads  int `alias:"Info.WeeklyDownloads"`
	MonthlyDownloads int `alias:"Info.MonthlyDownloads"`
	YearlyDownloads  int `alias:"Info.YearlyDownloads"`

	TotalStars   int `alias:"Info.TotalStars"`
	WeeklyStars  int `alias:"Info.WeeklyStars"`
	MonthlyStars int `alias:"Info.MonthlyStars"`
	YearlyStars  int `alias:"Info.YearlyStars"`
}

type ManyArts []struct {
	model.Arts

	Creator model.Users `alias:"Creator.*"`
	Cover   model.Files `alias:"Cover.*"`
	Tags    []model.Tags

	TotalDownloads   int `alias:"Info.TotalDownloads"`
	WeeklyDownloads  int `alias:"Info.WeeklyDownloads"`
	MonthlyDownloads int `alias:"Info.MonthlyDownloads"`
	YearlyDownloads  int `alias:"Info.YearlyDownloads"`

	TotalStars   int `alias:"Info.TotalStars"`
	WeeklyStars  int `alias:"Info.WeeklyStars"`
	MonthlyStars int `alias:"Info.MonthlyStars"`
	YearlyStars  int `alias:"Info.YearlyStars"`
}

type ManyArtsReq struct {
	Search string `form:"search"`
	Filter Filter
	Sort   Sort
	Page   int `form:"page"`
}

type Filter struct {
	Tags      []string `form:"filter.tags"`
	MinPrice  int      `form:"filter.minPrice"`
	MaxPrice  int      `form:"filter.maxPrice"`
	ImageExts []string `form:"filter.imageExts"`
}

type Sort struct {
	By  string `form:"sort.by"`
	Asc bool   `form:"sort.asc"`
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
