package tagsUsecase

import (
	"github.com/DeepAung/deep-art/modules/tags"
	"github.com/DeepAung/deep-art/modules/tags/tagsRepository"
)

type ITagsUsecase interface {
	GetTags() (*[]tags.Tag, error)
	CreateTag(req *tags.TagReq) error
	UpdateTag(req *tags.TagReq, id int) error
	DeleteTag(id int) error
}

type tagsUsecase struct {
	tagsRepository tagsRepository.ITagsRepository
}

func NewTagsUsecase(tagsRepository tagsRepository.ITagsRepository) ITagsUsecase {
	return &tagsUsecase{
		tagsRepository: tagsRepository,
	}
}

func (u *tagsUsecase) GetTags() (*[]tags.Tag, error) {
	return u.tagsRepository.GetTags()
}

func (u *tagsUsecase) CreateTag(req *tags.TagReq) error {
	return u.tagsRepository.CreateTag(req)
}

func (u *tagsUsecase) UpdateTag(req *tags.TagReq, id int) error {
	return u.tagsRepository.UpdateTag(req, id)
}

func (u *tagsUsecase) DeleteTag(id int) error {
	return u.tagsRepository.DeleteTag(id)
}
