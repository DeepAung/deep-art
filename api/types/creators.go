package types

type Creator struct {
	Id        int    `form:"id"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	AvatarURL string `form:"avatar_url"`
	Followers int    `form:"followers"`
}
