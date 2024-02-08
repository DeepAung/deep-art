package tags

type Tag struct {
	Id   int    `db:"id"   json:"id"   form:"id"`
	Name string `db:"name" json:"name" form:"name"`
}

type TagReq struct {
	Name string `db:"name" json:"name" form:"name"`
}
