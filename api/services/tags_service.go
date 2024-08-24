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

func (s *TagsSvc) GetTags() ([]model.Tags, error) {
	return s.tagsRepo.FindAllTags()
}

func (s *TagsSvc) CreateTag(name string) (model.Tags, error) {
	return s.tagsRepo.CreateTag(name)
}

func (s *TagsSvc) UpdateTag(id int, name string) (model.Tags, error) {
	if err := s.tagsRepo.UpdateTag(id, name); err != nil {
		return model.Tags{}, err
	}

	ID := int32(id)
	return model.Tags{
		ID:   &ID,
		Name: name,
	}, nil
}

func (s *TagsSvc) DeleteTag(id int) error {
	return s.tagsRepo.DeleteTag(id)
}
