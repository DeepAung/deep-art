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
	Email    string `form:"email"    validate:"required,email"`
	Password string `form:"password" validate:"required"`

	RedirectTo string `form:"redirect_to"`
}

type SetPasswordReq struct {
	Password        string `form:"password"         validate:"required"`
	ConfirmPassword string `form:"confirm_password" validate:"required"`
}

type SignUpDTO struct {
	SignUpReq

	RedirectTo string `form:"redirect_to"`
}

type SignUpReq struct {
	Username        string `form:"username"         validate:"required"`
	Email           string `form:"email"            validate:"required,email"`
	Password        string `form:"password"         validate:"required"`
	ConfirmPassword string `form:"confirm_password" validate:"required"`
	AvatarUrl       string `form:"avatar_url"`
}

type UpdateUserReq struct {
	Username  string `form:"username" validate:"required"`
	AvatarUrl string
}

type Token struct {
	Id           int    `form:"id"`
	AccessToken  string `form:"access_token"`
	RefreshToken string `form:"refresh_token"`
}

type OAuthInfo struct {
	ConnectGoogle bool
	ConnectGithub bool
}

type OAuthProvider struct {
	Provider string `alias:"Oauths.Provider"`
}

// Not used yet........

type ProviderEnum string

const (
	Google ProviderEnum = "google"
	Github ProviderEnum = "github"
)
