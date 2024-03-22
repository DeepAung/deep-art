package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserPassport struct {
	User  *User  `json:"user"  form:"user"`
	Token *Token `json:"token" form:"token"`
}

type User struct {
	Id        int    `db:"id"         json:"id"         form:"id"`
	Username  string `db:"username"   json:"username"   form:"username"`
	Email     string `db:"email"      json:"email"      form:"email"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url" form:"avatar_url"`
	IsAdmin   bool   `db:"is_admin"   json:"is_admin"   form:"is_admin"`
	Coin      int    `db:"coin"       json:"coin"       form:"coin"`
}

type UserWithPassword struct {
	Id        int    `db:"id"         json:"id"         form:"id"`
	Username  string `db:"username"   json:"username"   form:"username"`
	Email     string `db:"email"      json:"email"      form:"email"`
	Password  string `db:"password"   json:"password"   form:"password"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url" form:"avatar_url"`
	IsAdmin   bool   `db:"is_admin"   json:"is_admin"   form:"is_admin"`
	Coin      int    `db:"coin"       json:"coin"       form:"coin"`
}

type LoginReq struct {
	Email    string `db:"email"    json:"email"    form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type LogoutReq struct {
	TokenId int `db:"id" json:"token_id" form:"token_id"`
}

type RegisterReq struct {
	Username  string `db:"username"   json:"username"   form:"username"`
	Email     string `db:"email"      json:"email"      form:"email"`
	Password  string `db:"password"   json:"password"   form:"password"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url" form:"avatar_url"`
}

func (obj *RegisterReq) IsEmail() bool {
	matched, err := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(obj.Email))
	if err != nil {
		return false
	}
	return matched
}

func (obj *RegisterReq) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hash password failed: %v", err)
	}

	obj.Password = string(hashedPassword)
	return nil
}

type UpdateReq struct {
	Username  string `db:"username"   json:"username"   form:"username"`
	Password  string `db:"password"   json:"password"   form:"password"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url" form:"avatar_url"`
}

type Token struct {
	Id           int    `db:"id"            json:"id"            form:"id"`
	AccessToken  string `db:"access_token"  json:"access_token"  form:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token" form:"refresh_token"`
}

type RefreshTokensReq struct {
	RefreshToken string `db:"refresh_token" json:"refresh_token" form:"refresh_token"`
}

type TokenInfo struct {
	Id     int `db:"id"      json:"id"      form:"id"`
	UserId int `db:"user_id" json:"user_id" form:"user_id"`
}

type AdminTokenRes struct {
	AdminToken string `json:"admin_token" form:"admin_token"`
}

// type TokenReq struct {
// 	UserId       int    `db:"user_id"       json:"user_id"`
// 	AccessToken  string `db:"access_token"  json:"access_token"`
// 	RefreshToken string `db:"refresh_token" json:"refresh_token"`
// }

type OAuthReq struct {
	UserId int        `db:"user_id" json:"user_id" form:"user_id"`
	Social SocialEnum `db:"social"  json:"social"  form:"social"`
}

type OAuthCreateReq struct {
	UserId   int        `db:"user_id"   json:"user_id"   form:"user_id"`
	Social   SocialEnum `db:"social"    json:"social"    form:"social"`
	SocialId string     `db:"social_id" json:"social_id" form:"social_id"`
}

type OAuthDisconnectReq struct {
	Social SocialEnum `db:"social" json:"social" form:"social"`
}

type SocialEnum string

const (
	Google SocialEnum = "google"
	Github SocialEnum = "github"
)
