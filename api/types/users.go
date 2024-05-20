package types

type Passport struct {
	User  User  `form:"user"`
	Token Token `form:"token"`
}

type User struct {
	Id        int    `form:"id"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	AvatarUrl string `form:"avatar_url"`
	IsAdmin   bool   `form:"is_admin"`
	Coin      int    `form:"coin"`
}

type UserWithPassword struct {
	Id        int    `form:"id"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	Password  string `form:"password"`
	AvatarUrl string `form:"avatar_url"`
	IsAdmin   bool   `form:"is_admin"`
	Coin      int    `form:"coin"`
}

type SignInReq struct {
	Email    string `form:"email" validate:"required, email"`
	Password string `form:"password" validate:"required"`
}

type SignUpReq struct {
	Username  string `form:"username" validate:"required"`
	Email     string `form:"email" validate:"required, email"`
	Password  string `form:"password" validate:"required"`
	AvatarUrl string `form:"avatar_url"`
}

type SignOutReq struct {
	UserId  int `form:"user_id"`
	TokenId int `form:"token_id"`
}

type UpdateReq struct {
	Username  string `form:"username" validate:"required"`
	Email     string `form:"email" validate:"required, email"`
	AvatarUrl string `form:"avatar_url" validate:"required"`
}

type Token struct {
	Id           int    `form:"id"`
	AccessToken  string `form:"access_token"`
	RefreshToken string `form:"refresh_token"`
}

type RefreshTokensReq struct {
	RefreshToken string `form:"refresh_token" validate:"required"`
}

type TokenInfo struct {
	Id     int `form:"id"`
	UserId int `form:"user_id"`
}

type AdminTokenRes struct {
	AdminToken string `form:"admin_token"`
}

type OAuthReq struct {
	UserId int        `form:"user_id"`
	Social SocialEnum `form:"social"`
}

type OAuthCreateReq struct {
	UserId   int        `form:"user_id"`
	Social   SocialEnum `form:"social"`
	SocialId string     `form:"social_id"`
}

type OAuthDisconnectReq struct {
	Social SocialEnum `form:"social"`
}

type SocialEnum string

const (
	Google SocialEnum = "google"
	Github SocialEnum = "github"
)
