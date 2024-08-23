package utils

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	fb "path/filepath"
	"strings"
)

func Join(elem ...string) string {
	joined, _ := url.JoinPath(elem[0], elem[1:]...)
	return strings.ReplaceAll(joined, "\\", "/")
}

// url      = https://storage.googleapis.com/deep-art-bucket-dev/users/1/profile.jpg
// basePath = https://storage.googleapis.com/deep-art-bucket-dev
// bucket   = deep-art-bucket-dev
// dest     = users/1/profile.jpg
// dir      = users/1
// filename = profile.jpg
//
// url      = /static/storage/users/1/profile.jpg
// basePath = /static/storage
// dest     = users/1/profile.jpg
// dir      = users/1
// filename = profile.jpg
type UrlInfo struct {
	url      string
	basePath string
	dest     string
	dir      string
	filename string
}

func (u UrlInfo) Url() string      { return u.url }
func (u UrlInfo) BasePath() string { return u.basePath }
func (u UrlInfo) Dest() string     { return u.dest }
func (u UrlInfo) Dir() string      { return u.dir }
func (u UrlInfo) Filename() string { return u.filename }

func NewUrlInfoByDest(basePath string, dest string) UrlInfo {
	u := UrlInfo{
		basePath: basePath,
		dest:     dest,
	}

	idx := strings.LastIndex(u.dest, "/")
	u.filename = u.dest[idx+1:]

	url := Join(u.basePath, u.dest)
	u.url = strings.ReplaceAll(url, " ", "%20")

	u.dir = dest[0 : len(dest)-len(u.filename)]

	return u
}

func NewUrlInfoByURL(basePath string, url string) UrlInfo {
	u := UrlInfo{
		basePath: basePath,
		url:      url,
	}

	u.dest = u.url[len(u.basePath)+1:]

	idx := strings.LastIndex(u.dest, "/")
	u.filename = u.dest[idx+1:]

	u.dir = u.dest[0 : len(u.dest)-len(u.filename)]

	return u
}

func DownloadFiles(filespath, urls []string) (error, []string) {
	for i := range len(filespath) {
		err := DownloadFile(filespath[i], urls[i])
		if err != nil {
			return err, filespath[:i]
		}
	}

	return nil, filespath
}

func DeleteFiles(filespath []string) error {
	errorsMsg := make([]string, 0)
	for _, filepath := range filespath {
		if err := os.Remove(filepath); err != nil {
			errorsMsg = append(errorsMsg, err.Error())
		}
	}

	if len(errorsMsg) == 0 {
		return nil
	}

	return errors.New(strings.Join(errorsMsg, " "))
}

func DownloadFile(filepath, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func CreateZipFile(filespath []string, zipname string) (zipDest string, zipName string, err error) {
	zipDest = "tmp/" + zipname + ".zip"
	zipName = zipname + ".zip"

	var archive *os.File
	if archive, err = os.Create(zipDest); err != nil {
		return
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	for _, filepath := range filespath {
		var f *os.File
		if f, err = os.Open(filepath); err != nil {
			return
		}
		defer f.Close()

		var w io.Writer
		zipFilename := zipname + "/" + fb.Base(filepath)
		if w, err = zipWriter.Create(zipFilename); err != nil {
			return
		}

		if _, err = io.Copy(w, f); err != nil {
			return
		}
	}
	zipWriter.Close()

	return zipDest, zipName, nil
}
