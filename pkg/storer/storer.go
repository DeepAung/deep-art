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

type FileRes struct {
	Filename string `json:"filename" form:"filename"`
	Url      string `json:"url"      form:"url"`
}
