package services

import (
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/.gen/model"
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
)

var (
	ErrCodeExpired = httperror.New("code has already expired", http.StatusBadRequest)
	ErrCodeUsed    = httperror.New("code has already used", http.StatusBadRequest)
)

type CodesSvc struct {
	codesRepo *repositories.CodesRepo
	cfg       *config.Config
}

func NewCodesSvc(codesRepo *repositories.CodesRepo, cfg *config.Config) *CodesSvc {
	return &CodesSvc{
		codesRepo: codesRepo,
		cfg:       cfg,
	}
}

func (s *CodesSvc) UseCode(userId int, name string) error {
	// 1. check if code exist
	code, err := s.codesRepo.FindOneCodeByName(name)
	if err != nil {
		return err
	}

	// 2. check if code is expired
	if code.ExpTime.Before(time.Now()) {
		return ErrCodeExpired
	}

	// 3. check if code is used
	used, err := s.codesRepo.HasUsedCode(userId, int(*code.ID))
	if err != nil {
		return err
	}
	if used {
		return ErrCodeUsed
	}

	// 4. use the code
	return s.codesRepo.UseCode(userId, int(*code.ID))
}

func (s *CodesSvc) CreateCode(req types.CodeReq) (model.Codes, error) {
	return s.codesRepo.CreateCode(req)
}

func (s *CodesSvc) GetCodes() ([]model.Codes, error) {
	return s.codesRepo.FindAllCodes()
}

func (s *CodesSvc) UpdateCode(id int, req types.CodeReq) (model.Codes, error) {
	if err := s.codesRepo.UpdateCode(id, req); err != nil {
		return model.Codes{}, err
	}

	ID := int32(id)
	return model.Codes{
		ID:      &ID,
		Name:    req.Name,
		Value:   int32(req.Value),
		ExpTime: req.ExpTime.Time,
	}, nil
}

func (s *CodesSvc) DeleteCode(id int) error {
	return s.codesRepo.DeleteCode(id)
}
