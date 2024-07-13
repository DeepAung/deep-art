package services

import (
	"github.com/DeepAung/deep-art/.gen/model"
	"github.com/DeepAung/deep-art/api/repositories"
)

type TagsSvc struct {
	tagsRepo *repositories.TagsRepo
}

func NewTagsSvc(tagsRepo *repositories.TagsRepo) *TagsSvc {
	return &TagsSvc{
		tagsRepo: tagsRepo,
	}
}

func (s *TagsSvc) FindAllTags() ([]model.Tags, error) {
	return s.tagsRepo.FindAllTags()
}
