package mystorage

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/pkg/utils"
)

type IStorage interface {
	UploadFiles(files []*multipart.FileHeader, dir string) ([]*FileRes, error)
	DeleteFiles(destinations []string) error
}

type FileRes struct {
	Filename string `json:"filename" form:"filename"`
	Url      string `json:"url"      form:"url"`
}

type fileInfo struct {
	File        *multipart.FileHeader
	Dir         string // "images"
	Filename    string // "example.png"
	Extension   string // "png"
	Destination string // "images/example.png"
}

func newFilesInfo(files []*multipart.FileHeader, dir string) []*fileInfo {
	filesInfo := make([]*fileInfo, len(files))

	for i, file := range files {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		filename := utils.RandFilename(ext)

		filesInfo[i] = &fileInfo{
			File:        file,
			Dir:         dir,
			Filename:    filename,
			Extension:   ext,
			Destination: utils.Join(dir, filename),
		}
	}

	return filesInfo
}

func newGCPFileRes(fileName, destination, bucket string) *FileRes {
	return &FileRes{
		Filename: fileName,
		Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, destination),
	}
}

func newLocalFileRes(cfg config.IAppConfig, fileName, destination string) *FileRes {
	return &FileRes{
		Filename: fileName,
		Url:      fmt.Sprintf("http://%s:%d/%s", cfg.Host(), cfg.Port(), destination),
	}
}
