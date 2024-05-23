package utils

import (
	"path/filepath"
	"strings"
)

func Join(elem ...string) string {
	return strings.ReplaceAll(filepath.Join(elem...), "\\", "/")
}

// url      = https://storage.googleapis.com/deep-art-bucket-dev/users/1/profile.jpg
// website  = https://storage.googleapis.com
// bucket   = deep-art-bucket-dev
// dest     = users/1/profile.jpg
// filename = profile.jpg
type GcpUrl struct {
	website string
	bucket  string
	dest    string
}

func NewGcpUrl(website string, bucket string, dest string) GcpUrl {
	return GcpUrl{
		website: website,
		bucket:  bucket,
		dest:    dest,
	}
}

func (u GcpUrl) Filename() string {
	idx := strings.LastIndex(u.dest, "/")
	return u.dest[idx+1:]
}

func (u GcpUrl) Url() string {
	url := Join(u.website, u.bucket, u.dest)
	return strings.ReplaceAll(url, " ", "%20")
}

type LocalUrl struct {
	basePath string
	dest     string
}

// url      = ./static/storage/users/1/profile.jpg
// basePath  = ./static/storage
// dest     = users/1/profile.jpg
// filename = profile.jpg
func NewLocalUrl(basePath string, dest string) LocalUrl {
	return LocalUrl{
		basePath: basePath,
		dest:     dest,
	}
}

func (u LocalUrl) Filename() string {
	idx := strings.LastIndex(u.dest, "/")
	return u.dest[idx+1:]
}

func (u LocalUrl) Url() string {
	url := Join(u.basePath, u.dest)
	return strings.ReplaceAll(url, " ", "%20")
}
