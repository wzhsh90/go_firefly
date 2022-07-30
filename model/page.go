package models

type PageModel struct {
	PageIndex  int         `json:"-"`
	PageSize   int         `json:"pageSize"`
	Pages      int         `json:"pages"`
	Page       int         `json:"page"`
	ItemsCount int         `json:"itemsCount"`
	Data       interface{} `json:"data"`
}
type PageModelLay struct {
	Code      int         `json:"code"`
	PageIndex int         `json:"-"`
	PageSize  int         `json:"pageSize"`
	Pages     int         `json:"pages"`
	Page      int         `json:"page"`
	Records   int         `json:"records"`
	Rows      interface{} `json:"rows"`
}

func maxFn(x, y int) int {
	if x < y {
		return y
	}
	return x
}
func (m *PageModel) BuildPageInfo(pageNo, pageSize, records int) {
	m.Page = maxFn(pageNo, 1)
	m.PageSize = pageSize
	m.Pages = (records + pageSize - 1) / pageSize
	if pageSize >= 2 && records <= pageSize*(pageNo-1) {
		m.Page = m.Pages
	}
	m.Page = maxFn(m.Page, 1)
	var pageIndex = (m.Page - 1) * m.PageSize
	m.ItemsCount = records
	m.PageIndex = pageIndex
}

func (m *PageModelLay) BuildPageInfo(pageNo, pageSize, records int) {
	m.Code = 0
	m.Page = maxFn(pageNo, 1)
	m.PageSize = pageSize
	m.Pages = (records + pageSize - 1) / pageSize
	if pageSize >= 2 && records <= pageSize*(pageNo-1) {
		m.Page = m.Pages
	}
	m.Page = maxFn(m.Page, 1)
	var pageIndex = (m.Page - 1) * m.PageSize
	m.Records = records
	m.PageIndex = pageIndex
}
