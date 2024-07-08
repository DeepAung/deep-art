package storer

import "mime/multipart"

type Storer interface {
	UploadFiles(files []*multipart.FileHeader, dir string) ([]FileRes, error)
	DeleteFiles(destinations []string) error
}

type FileReq struct {
	file *multipart.FileHeader
	dir  string
}

type FileRes interface {
	Url() string
	BasePath() string
	Dest() string
	Dir() string
	Filename() string
}
