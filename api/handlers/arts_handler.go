package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/components"
	"github.com/labstack/echo/v4"
)

type ArtsHandler struct {
	artsSvc *services.ArtsSvc
	cfg     *config.Config
}

func NewArtsHandler(artsSvc *services.ArtsSvc, cfg *config.Config) *ArtsHandler {
	return &ArtsHandler{
		artsSvc: artsSvc,
		cfg:     cfg,
	}
}

func (h *ArtsHandler) CreateArt(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(c, components.Error, ErrPayloadNotFound)
	}

	var dto types.FullArtDTO
	if err := c.Bind(&dto); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&dto); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	covers, ok := form.File["cover"]
	if !ok || len(covers) == 0 {
		return utils.Render(c, components.Error("no \"cover\" field"), http.StatusBadRequest)
	}
	dto.Cover = covers[0]
	dto.Files, ok = form.File["files"]
	if !ok || len(dto.Files) == 0 {
		return utils.Render(c, components.Error("no \"files\" field"), http.StatusBadRequest)
	}

	if err := h.artsSvc.CreateArt(payload.UserId, dto); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Redirect", "/creator")
	return nil
}

func (h *ArtsHandler) UpdateArt(c echo.Context) error {
	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	var dto types.ArtDTO
	if err := c.Bind(&dto); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&dto); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	if err := h.artsSvc.UpdateArtInfo(types.UpdateArtInfoReq{
		ArtId:       artId,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		TagsID:      dto.TagsID,
	}); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (h *ArtsHandler) DeleteArt(c echo.Context) error {
	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	if err := h.artsSvc.DeleteArt(artId); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Redirect", "/creator")
	return c.NoContent(http.StatusOK)
}

func (h *ArtsHandler) DownloadArt(c echo.Context) error {
	renderWentWrongError := func() error {
		return utils.Render(
			c,
			components.Error("something went wrong. please try again."),
			http.StatusInternalServerError,
		)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	art, err := h.artsSvc.FindOneArt(artId)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	dests := make([]string, len(art.Files))
	urls := make([]string, len(art.Files))
	for i, file := range art.Files {
		dests[i] = "tmp/" + file.Filename
		urls[i] = file.URL
		if h.cfg.App.StorerType == config.LocalType {
			urls[i] = "http://" + c.Request().Host + urls[i]
		}
	}

	err, filespath := utils.DownloadFiles(dests, urls)
	if err != nil {
		tryDeleteFiles(err, filespath) // Rollback stuff
		return renderWentWrongError()
	}

	var zipDest, zipName string
	if zipDest, zipName, err = utils.CreateZipFile(dests, art.Name); err != nil {
		tryDeleteFiles(err, append(filespath, zipDest))
		return renderWentWrongError()
	}

	if err := c.Attachment(zipDest, zipName); err != nil {
		tryDeleteFiles(err, append(filespath, zipDest))
		return renderWentWrongError()
	}

	return utils.DeleteFiles(append(filespath, zipDest))
}

func tryDeleteFiles(causeErr error, files []string) {
	if err := utils.DeleteFiles(files); err != nil {
		slog.Error(causeErr.Error() + " " + err.Error())
	} else {
		slog.Error(causeErr.Error())
	}
}

func (h *ArtsHandler) UploadFiles(c echo.Context) error {
	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	files, ok := form.File["files"]
	if !ok || len(files) == 0 {
		return utils.Render(c, components.Error("no \"files\" field"), http.StatusBadRequest)
	}

	if err := h.artsSvc.UploadFiles(artId, files); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (h *ArtsHandler) DeleteFile(c echo.Context) error {
	artId, err := strconv.Atoi(c.Param("artId"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	fileId, err := strconv.Atoi(c.Param("fileId"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	if err := h.artsSvc.DeleteFile(artId, fileId); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (h *ArtsHandler) ReplaceCover(c echo.Context) error {
	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	covers, ok := form.File["cover"]
	if !ok || len(covers) == 0 {
		return utils.Render(c, components.Error("no \"cover\" field"), http.StatusBadRequest)
	}

	if err := h.artsSvc.ReplaceCover(artId, covers[0]); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (h *ArtsHandler) FindManyArts(c echo.Context) error {
	req, err := h.getValidatedManyArtsReq(c)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}
	if err := utils.Validate(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	arts, err := h.artsSvc.FindManyArts(req)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	time.Sleep(300 * time.Millisecond)
	return utils.Render(c, components.ManyArts(arts, false), http.StatusOK)
}

func (h *ArtsHandler) FindManyArtsWithArtType(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(c, components.Error, ErrPayloadNotFound)
	}

	req, err := h.getValidatedManyArtsReq(c)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	artType := c.QueryParam("artType")
	withEdit := false
	if c.QueryParam("withEdit") == "true" {
		withEdit = true
	}

	var arts types.ManyArtsRes
	switch artType {
	case "starred":
		arts, err = h.artsSvc.FindManyStarredArts(payload.UserId, req)
	case "bought":
		arts, err = h.artsSvc.FindManyBoughtArts(payload.UserId, req)
	case "created":
		arts, err = h.artsSvc.FindManyCreatedArts(payload.UserId, req)
	default:
		return utils.Render(c, components.Error("invalid arts type"), http.StatusBadRequest)
	}
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	time.Sleep(300 * time.Millisecond)
	return utils.Render(c, components.ManyArts(arts, withEdit), http.StatusOK)
}

func (h *ArtsHandler) getValidatedManyArtsReq(c echo.Context) (types.ManyArtsReq, error) {
	var dto types.ManyArtsDTO
	if err := c.Bind(&dto); err != nil {
		return types.ManyArtsReq{}, httperror.New(err.Error(), http.StatusBadRequest)
	}

	var (
		filter     types.Filter
		sort       types.Sort
		pagination types.Pagination
	)

	err := json.Unmarshal([]byte(dto.Filter), &filter)
	if err != nil {
		return types.ManyArtsReq{}, httperror.New(
			"invalid filter body request",
			http.StatusBadRequest,
		)
	}

	err = json.Unmarshal([]byte(dto.Sort), &sort)
	if err != nil {
		return types.ManyArtsReq{}, httperror.New(
			"invalid sort body request",
			http.StatusBadRequest,
		)
	}

	err = json.Unmarshal([]byte(dto.Pagination), &pagination)
	if err != nil {
		return types.ManyArtsReq{}, httperror.New(
			"invalid pagination body request",
			http.StatusBadRequest,
		)
	}

	req := types.ManyArtsReq{
		Search:     dto.Search,
		Filter:     filter,
		Sort:       sort,
		Pagination: pagination,
	}

	if err := utils.Validate(&req); err != nil {
		return types.ManyArtsReq{}, httperror.New(err.Error(), http.StatusBadRequest)
	}

	return req, nil
}

func (h *ArtsHandler) BuyArt(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(c, components.Error, ErrPayloadNotFound)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.ErrToast, err)
	}

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return utils.RenderError(c, components.ErrToast, err)
	}

	if err := h.artsSvc.BuyArt(payload.UserId, artId, price); err != nil {
		return utils.RenderError(c, components.ErrToast, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return nil
}

func (h *ArtsHandler) ToggleStar(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(
			c,
			components.Error,
			ErrPayloadNotFound,
		)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	isStarred, err := h.artsSvc.ToggleStar(payload.UserId, artId)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.StarButton(artId, isStarred), http.StatusOK)
}
