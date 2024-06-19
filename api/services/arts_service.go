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

func (s *ArtsSvc) FindManyArts(req types.ManyArtsReq) (types.ManyArtsRes, error) {
	return s.artsRepo.FindManyArts(req)
}

func (s *ArtsSvc) FindOneArt(id int) (types.Art, error) {
	return s.artsRepo.FindOneArt(id)
}

func (s *ArtsSvc) ToggleStar(userId, artId int) (bool, error) {
	isStarred, err := s.IsStarred(userId, artId)
	if err != nil {
		return false, err
	}

	if isStarred {
		err = s.UnStar(userId, artId)
	} else {
		err = s.Star(userId, artId)
	}

	if err != nil {
		return false, err
	}

	return !isStarred, nil
}

func (s *ArtsSvc) IsStarred(userId, artId int) (bool, error) {
	return s.artsRepo.HasUsersStarredArts(userId, artId)
}

func (s *ArtsSvc) Star(userId, artId int) error {
	return s.artsRepo.CreateUsersStarredArts(userId, artId)
}

func (s *ArtsSvc) UnStar(userId, artId int) error {
	return s.artsRepo.DeleteUsersStarredArts(userId, artId)
}
