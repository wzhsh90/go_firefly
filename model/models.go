package models

type LoginUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	Id      string `form:"id" json:"id"`
	ComDesc string `form:"com_desc" json:"com_desc"`
	ComName string `form:"com_name" json:"com_name"`
}
