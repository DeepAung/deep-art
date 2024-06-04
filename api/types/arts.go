package types

import "github.com/DeepAung/deep-art/.gen/model"

type ManyArts []struct {
	model.Arts

	Creator model.Users `alias:"Creator.*"`
	Cover   model.Files `alias:"Cover.*"`
	Files   []model.Files
	Tags    []model.Tags

	TotalDownloads int `alias:"Info.TotalDownloads"`
	TotalStars     int `alias:"Info.TotalStars"`
}

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

type ManyArtsReq struct {
	Page int `form:"page"`
}
