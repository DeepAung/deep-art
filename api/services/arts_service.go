package services

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/DeepAung/deep-art/pkg/utils"
)

var (
// ErrInvalidEmailOrPassword = httperror.New("invalid email or password", http.StatusBadRequest)
// ErrInvalidRefreshToken    = httperror.New("invalid refresh token", http.StatusBadRequest)
)

type ArtsSvc struct {
	artsRepo *repositories.ArtsRepo
	storer   storer.Storer
	cfg      *config.Config
}

func NewArtsSvc(
	artsRepo *repositories.ArtsRepo,
	storer storer.Storer,
	cfg *config.Config,
) *ArtsSvc {
	return &ArtsSvc{
		artsRepo: artsRepo,
		storer:   storer,
		cfg:      cfg,
	}
}

// create art
// get art id
// upload cover & files
// update art coverURL & filesURL
func (s *ArtsSvc) CreateArt(creatorId int, dto types.ArtDTO) error {
	ctx, cancel, tx, err := s.artsRepo.BeginTx()
	defer cancel()
	if err != nil {
		return err
	}

	// create art
	createReq := types.CreateArtReq{
		ArtDTO:    dto,
		CreatorId: creatorId,
	}
	artId, err := s.artsRepo.CreateArtWithDB(ctx, tx, createReq)
	if err != nil {
		return err
	}

	// upload cover
	coverDir := fmt.Sprint("/arts/cover/", artId)
	coverRes, err := s.storer.UploadFiles([]*multipart.FileHeader{dto.Cover}, coverDir)
	if err != nil {
		return err
	}
	coverDest := []string{coverRes[0].Dest()}

	// upload files
	filesDir := fmt.Sprint("/arts/files/", artId)
	filesRes, err := s.storer.UploadFiles(dto.Files, filesDir)
	if err != nil {
		_ = s.storer.DeleteFiles(coverDest) // rollback process
		return err
	}
	filesDest := utils.Map(filesRes, func(t storer.FileRes) string { return t.Dest() })

	coverURL := coverRes[0].Url()
	filesURL := utils.Map(filesRes, func(t storer.FileRes) string { return t.Url() })

	// update art
	updateReq := types.UpdateArtReq{
		ArtDTO:   dto,
		ArtID:    artId,
		CoverURL: coverURL,
		FilesURL: filesURL,
	}
	if err := s.artsRepo.UpdateArtWithDB(ctx, tx, updateReq); err != nil {
		_ = s.storer.DeleteFiles(append(filesDest, coverDest...)) // rollback process
		return err
	}

	return tx.Commit()
}

func (s *ArtsSvc) FindManyArts(req types.ManyArtsReq) (types.ManyArtsRes, error) {
	return s.artsRepo.FindManyArts(req)
}

func (s *ArtsSvc) FindManyStarredArts(
	userId int,
	req types.ManyArtsReq,
) (types.ManyArtsRes, error) {
	return s.artsRepo.FindManyStarredArts(userId, req)
}

func (s *ArtsSvc) FindManyBoughtArts(
	userId int,
	req types.ManyArtsReq,
) (types.ManyArtsRes, error) {
	return s.artsRepo.FindManyBoughtArts(userId, req)
}

func (s *ArtsSvc) FindManyCreatedArts(
	userId int,
	req types.ManyArtsReq,
) (types.ManyArtsRes, error) {
	return s.artsRepo.FindManyCreatedArts(userId, req)
}

func (s *ArtsSvc) FindOneArt(id int) (types.Art, error) {
	return s.artsRepo.FindOneArt(id)
}

func (s *ArtsSvc) BuyArt(userId, artId, price int) error {
	bought, err := s.artsRepo.HasUsersBoughtArts(userId, artId)
	if err != nil {
		return err
	}
	if bought {
		return httperror.New("user already bought this art", http.StatusBadRequest)
	}

	coin, err := s.artsRepo.FindUserCoin(userId)
	if err != nil {
		return err
	}
	if coin < price {
		return httperror.New("not enough coin to buy this art", http.StatusBadRequest)
	}

	return s.artsRepo.BuyArt(userId, artId, price)
}

func (s *ArtsSvc) IsBought(userId, artId int) (bool, error) {
	return s.artsRepo.HasUsersBoughtArts(userId, artId)
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
