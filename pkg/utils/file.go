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
	Url      string
	BasePath string
	Dest     string
	Dir      string
	Filename string
}

func NewUrlInfoByDest(basePath string, dest string) UrlInfo {
	u := UrlInfo{
		BasePath: basePath,
		Dest:     dest,
	}

	idx := strings.LastIndex(u.Dest, "/")
	u.Filename = u.Dest[idx+1:]

	url := Join(u.BasePath, u.Dest)
	u.Url = strings.ReplaceAll(url, " ", "%20")

	u.Dir = dest[0 : len(dest)-len(u.Filename)]

	return u
}

func NewUrlInfoByURL(basePath string, url string) UrlInfo {
	u := UrlInfo{
		BasePath: basePath,
		Url:      url,
	}

	u.Dest = u.Url[len(u.BasePath)+1:]

	idx := strings.LastIndex(u.Dest, "/")
	u.Filename = u.Dest[idx+1:]

	u.Dir = u.Dest[0 : len(u.Dest)-len(u.Filename)]

	return u
}
