package models

type LoginUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	Id      string `gorm:"id" form:"id" json:"id"`
	ComDesc string `gorm:"com_desc" form:"com_desc" json:"com_desc"`
	ComName string `gorm:"com_name" form:"com_name" json:"com_name"`
	Flag    int64  `gorm:"flag" form:"flag" json:"flag"`
	Age     int64  `gorm:"age" form:"age" json:"age"`
}
