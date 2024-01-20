package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserPassport struct {
	User  *User  `json:"user"`
	Token *Token `json:"token"`
}

type User struct {
	Id        int    `db:"id"         json:"id"`
	Username  string `db:"username"   json:"username"`
	Email     string `db:"email"      json:"email"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url"`
}

type UserWithPassword struct {
	Id        int    `db:"id"         json:"id"`
	Username  string `db:"username"   json:"username"`
	Email     string `db:"email"      json:"email"`
	Password  string `db:"password"   json:"password"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url"`
}

type LoginReq struct {
	Email    string `db:"email"    json:"email"`
	Password string `db:"password" json:"password"`
}

type RegisterReq struct {
	Username  string `db:"username"   json:"username"`
	Email     string `db:"email"      json:"email"`
	Password  string `db:"password"   json:"password"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url"`
}

type RegisterOAuthReq struct {
	Username  string `db:"username"   json:"username"`
	Email     string `db:"email"      json:"email"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url"`
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
	Username  string `db:"username"   json:"username"`
	Password  string `db:"password"   json:"password"`
	AvatarUrl string `db:"avatar_url" json:"avatar_url"`
}

type Token struct {
	Id           int    `db:"id"            json:"id"`
	AccessToken  string `db:"access_token"  json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type TokenReq struct {
	UserId       int    `db:"user_id"       json:"user_id"`
	AccessToken  string `db:"access_token"  json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type OAuth struct {
	Id       int        `db:"id"        json:"id"`
	UserId   int        `db:"user_id"   json:"user_id"`
	Social   SocialEnum `db:"social"    json:"social"`
	SocialId string     `db:"social_id" json:"social_id"`
	// CreatedAt time.Time  `db:"created_at" json:"created_at"`
	// UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type OAuthReq struct {
	UserId   int        `db:"user_id"   json:"user_id"`
	Social   SocialEnum `db:"social"    json:"social"`
	SocialId string     `db:"social_id" json:"social_id"`
}

type SocialEnum string

const (
	Google SocialEnum = "google"
	Github SocialEnum = "github"
)
