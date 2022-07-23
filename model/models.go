package models

import "github.com/upper/db/v4"

type LoginUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	Id      string `db:"id" form:"id" json:"id"`
	ComDesc string `db:"com_desc" form:"com_desc" json:"com_desc"`
	ComName string `db:"com_name" form:"com_name" json:"com_name"`
}

func (b *Company) Store(sess db.Session) db.Store {
	return sess.Collection("sys_company_t")
}
