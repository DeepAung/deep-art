package services

import "github.com/DeepAung/deep-art/api/repositories"

type TagsSvc struct {
	usersRepo *repositories.TagsRepo
}

func NewTagsSvc(usersRepo *repositories.TagsRepo) *TagsSvc {
	return &TagsSvc{
		usersRepo: usersRepo,
	}
}
