package storer

import (
	"io"
)

type Storer interface {
	UploadFile(file io.Reader, dest string) (FileRes, error)
	DeleteFile(dest string) error
	UploadFiles(files []io.Reader, dests []string) ([]FileRes, error)
	DeleteFiles(dests []string) error
}

type FileRes interface {
	Url() string
	BasePath() string
	Dest() string
	Dir() string
	Filename() string
}
