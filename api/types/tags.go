package types

type TagReq struct {
	Name string `form:"name" validate:"required"`
}
