package types

type CodeReq struct {
	Name    string     `form:"name"    validate:"required"`
	Value   int        `form:"value"   validate:"required,gte=0"`
	ExpTime CustomTime `form:"expTime" validate:"required"`
}
