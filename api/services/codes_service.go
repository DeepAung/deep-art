package services

import (
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/api/repositories"
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

// 1. check if code exist
// 2. check if code is expired
// 3. check if code is used
// 4. use the code
func (s *CodesSvc) UseCode(userId int, name string) error {
	code, err := s.codesRepo.FindOneCodeByName(name)
	if err != nil {
		return err
	}

	if code.ExpTime.Before(time.Now()) {
		return ErrCodeExpired
	}

	used, err := s.codesRepo.HasUsedCode(userId, int(*code.ID))
	if err != nil {
		return err
	}
	if used {
		return ErrCodeUsed
	}

	return s.codesRepo.UseCode(userId, int(*code.ID))
}
