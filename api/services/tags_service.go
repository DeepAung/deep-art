package services

import (
	"github.com/DeepAung/deep-art/.gen/model"
	"github.com/DeepAung/deep-art/api/repositories"
)

type TagsSvc struct {
	usersRepo *repositories.TagsRepo
}

func NewTagsSvc(usersRepo *repositories.TagsRepo) *TagsSvc {
	return &TagsSvc{
		usersRepo: usersRepo,
	}
}

func (s *TagsSvc) FindAllTags() ([]model.Tags, error) {
	return s.usersRepo.FindAllTags()
}
