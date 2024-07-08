package utils_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/DeepAung/deep-art/pkg/utils"
)

var answers = []utils.UrlInfo{
	{
		Ur:       "/static/storage/users/1/profile.jpg",
		BasePath: "/static/storage",
		Dest:     "users/1/profile.jpg",
		Dir:      "users/1/",
		Filename: "profile.jpg",
	},
	{
		Ur:       "https://storage.googleapis.com/deep-art-bucket-dev/users/1/profile.jpg",
		BasePath: "https://storage.googleapis.com/deep-art-bucket-dev",
		Dest:     "users/1/profile.jpg",
		Dir:      "users/1/",
		Filename: "profile.jpg",
	},
}

func TestNewUrlInfoByDest(t *testing.T) {
	for i, ans := range answers {
		basePath := ans.BasePath
		dest := ans.Dest
		u := utils.NewUrlInfoByDest(basePath, dest)

		checkEqual(i, u, ans)
	}
}

func TestNewUrlInfoByUrl(t *testing.T) {
	for i, ans := range answers {
		basePath := ans.BasePath
		url := ans.Ur
		u := utils.NewUrlInfoByURL(basePath, url)

		checkEqual(i, u, ans)
	}
}

func checkEqual(idx int, u, ans any) {
	reU := reflect.ValueOf(u)
	reAns := reflect.ValueOf(ans)

	for j := 0; j < reU.NumField(); j++ {
		if !reflect.DeepEqual(reU.Field(j).Interface(), reAns.Field(j).Interface()) {
			log.Fatalf(
				"idx %d: field %s not equal | %v | %v",
				idx,
				reU.Type().Field(j).Name,
				reU.Field(j).Interface(),
				reAns.Field(j).Interface(),
			)
		}
	}
}
