package services

import (
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
)

var (
// ErrInvalidEmailOrPassword = httperror.New("invalid email or password", http.StatusBadRequest)
// ErrInvalidRefreshToken    = httperror.New("invalid refresh token", http.StatusBadRequest)
)

type ArtsSvc struct {
	artsRepo *repositories.ArtsRepo
	cfg      *config.Config
}

func NewArtsSvc(artsRepo *repositories.ArtsRepo, cfg *config.Config) *ArtsSvc {
	return &ArtsSvc{
		artsRepo: artsRepo,
		cfg:      cfg,
	}
}

func (s *ArtsSvc) FindManyArts(req types.ManyArtsReq) (types.ManyArts, error) {
	return s.artsRepo.FindManyArts(req)
}

func (s *ArtsSvc) FindOneArt(id int) (types.Art, error) {
	return s.artsRepo.FindOneArt(id)
}
