package mytoken_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/DeepAung/deep-art/pkg/mytoken"
)

func TestGenAndParseToken(t *testing.T) {
	secret := []byte("mysecret")

	payload := mytoken.Payload{
		UserId:   1111,
		Username: "test",
	}
	tokenString, err := mytoken.GenerateToken(mytoken.Access, 1000*time.Second, secret, payload)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := mytoken.ParseToken(mytoken.Access, secret, tokenString)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(payload, claims.Payload) {
		t.Fatal("not the same payload")
	}

}
