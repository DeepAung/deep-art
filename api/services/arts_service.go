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
func (s *ArtsSvc) CreateArt(creatorId int, dto types.FullArtDTO) error {
	ctx, cancel, tx, err := s.artsRepo.BeginTx()
	defer cancel()
	if err != nil {
		return err
	}

	// create art
	createReq := types.CreateArtReq{
		CreatorId:   creatorId,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		TagsID:      dto.TagsID,
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
	filesname := utils.Map(filesRes, func(t storer.FileRes) string { return t.Filename() })

	// update art
	updateReq := types.UpdateArtFilesReq{
		ArtId:     artId,
		CoverURL:  coverURL,
		FilesURL:  filesURL,
		FilesName: filesname,
	}
	if err := s.artsRepo.UpdateArtCoverAndFilesWithDB(ctx, tx, updateReq); err != nil {
		_ = s.storer.DeleteFiles(append(filesDest, coverDest...)) // rollback process
		return err
	}

	return tx.Commit()
}

func (s *ArtsSvc) UpdateArtInfo(req types.UpdateArtInfoReq) error {
	return s.artsRepo.UpdateArtInfo(req)
}

func (s *ArtsSvc) DeleteArt(artId int) error {
	ctx, cancel, tx, err := s.artsRepo.BeginTx()
	defer cancel()
	if err != nil {
		return err
	}

	// delete art
	if err := s.artsRepo.DeleteArtWithDB(ctx, tx, artId); err != nil {
		return err
	}

	// delete files
	files, err := s.artsRepo.FindManyFilesByArtId(artId)
	if err != nil {
		return err
	}

	var filesDest []string
	for _, file := range files {
		fileInfo := utils.NewUrlInfoByURL(s.cfg.App.BasePath, file.URL)
		filesDest = append(filesDest, fileInfo.Dest())
	}
	if err := s.storer.DeleteFiles(filesDest); err != nil {
		return err
	}

	// delete cover
	coverURL, err := s.artsRepo.FindOneCoverURL(artId)
	if err != nil {
		return err
	}
	coverInfo := utils.NewUrlInfoByURL(s.cfg.App.BasePath, coverURL)
	if err := s.storer.DeleteFiles([]string{coverInfo.Dest()}); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ArtsSvc) UploadFiles(artId int, files []*multipart.FileHeader) error {
	ctx, cancel, tx, err := s.artsRepo.BeginTx()
	defer cancel()
	if err != nil {
		return err
	}

	filesDir := fmt.Sprint("/arts/files/", artId)
	filesInfo := utils.Map(files, func(file *multipart.FileHeader) storer.FileRes {
		return utils.NewUrlInfoByDest(s.cfg.App.BasePath, filesDir+"/"+file.Filename)
	})
	filesURL := utils.Map(filesInfo, func(file storer.FileRes) string { return file.Url() })
	filesName := utils.Map(filesInfo, func(file storer.FileRes) string { return file.Filename() })

	fmt.Println("artId: ", artId)
	fmt.Println("filesURL: ", filesURL)
	fmt.Println("filesName: ", filesName)

	if err := s.artsRepo.InsertArtFilesWithDB(ctx, tx, artId, filesURL, filesName); err != nil {
		return err
	}

	if _, err := s.storer.UploadFiles(files, filesDir); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ArtsSvc) DeleteFile(artId, fileId int) error {
	ctx, cancel, tx, err := s.artsRepo.BeginTx()
	defer cancel()
	if err != nil {
		return err
	}

	if err := s.artsRepo.DeleteArtFilesWithDB(ctx, tx, fileId); err != nil {
		return err
	}

	file, err := s.artsRepo.FindOneFile(fileId)
	if err != nil {
		return err
	}

	fileInfo := utils.NewUrlInfoByURL(s.cfg.App.BasePath, file.URL)
	if err := s.storer.DeleteFiles([]string{fileInfo.Dest()}); err != nil {
		// there is no image but we want to delete it. so just do nothing
		return tx.Commit()
	}

	return tx.Commit()
}

func (s *ArtsSvc) ReplaceCover(artId int, cover *multipart.FileHeader) error {
	ctx, cancel, tx, err := s.artsRepo.BeginTx()
	defer cancel()
	if err != nil {
		return err
	}

	oldCoverURL, err := s.artsRepo.FindOneCoverURL(artId)
	if err != nil {
		return err
	}
	oldCoverInfo := utils.NewUrlInfoByURL(s.cfg.App.BasePath, oldCoverURL)

	newCoverDest := fmt.Sprintf("/arts/cover/%d/%s", artId, cover.Filename)
	newCoverInfo := utils.NewUrlInfoByDest(s.cfg.App.BasePath, newCoverDest)

	if err := s.artsRepo.UpdateArtCoverWithDB(ctx, tx, artId, newCoverInfo.Url()); err != nil {
		return err
	}

	if err := s.storer.DeleteFiles([]string{oldCoverInfo.Dest()}); err != nil {
		return err
	}
	if _, err := s.storer.UploadFiles([]*multipart.FileHeader{cover}, newCoverInfo.Dir()); err != nil {
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

func (s *ArtsSvc) Owned(userId, artId int) (bool, error) {
	creatorId, err := s.artsRepo.FindCreatorID(artId)
	if err != nil {
		return false, err
	}

	return creatorId == userId, nil
}
