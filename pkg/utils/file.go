package utils

import (
	"net/url"
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
