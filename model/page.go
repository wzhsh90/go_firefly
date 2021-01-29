package models

type PageModel struct {
	PageIndex  int64
	PageSize   int64       `json:"pageSize"`
	Pages      int64       `json:"pages"`
	Page       int64       `json:"page"`
	ItemsCount int64       `json:"itemsCount"`
	Data       interface{} `json:"data"`
}
type PageModelLay struct {
	Code      int `json:"code"`
	PageIndex int64
	PageSize  int64       `json:"pageSize"`
	Pages     int64       `json:"pages"`
	Page      int64       `json:"page"`
	Records   int64       `json:"records"`
	Rows      interface{} `json:"rows"`
}

func maxFn(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
func (m *PageModel) BuildPageInfo(pageNo, pageSize, records int64) {
	m.Page = maxFn(pageNo, 1)
	if pageSize == -1 {
		m.PageSize = records
		m.Pages = 1
	} else {
		m.PageSize = pageSize
		m.Pages = (records + pageSize - 1) / pageSize
		if pageSize >= 2 && records <= pageSize*(pageNo-1) {
			m.Page = m.Pages
		}
	}
	m.Page = maxFn(m.Page, 1)
	var pageIndex = (m.Page - 1) * m.PageSize
	m.ItemsCount = records
	m.PageIndex = pageIndex
}

func (m *PageModelLay) BuildPageInfo(pageNo, pageSize, records int64) {
	m.Code = 0
	m.Page = maxFn(pageNo, 1)
	if pageSize == -1 {
		m.PageSize = records
		m.Pages = 1
	} else {
		m.PageSize = pageSize
		m.Pages = (records + pageSize - 1) / pageSize
		if pageSize >= 2 && records <= pageSize*(pageNo-1) {
			m.Page = m.Pages
		}
	}
	m.Page = maxFn(m.Page, 1)
	var pageIndex = (m.Page - 1) * m.PageSize
	m.Records = records
	m.PageIndex = pageIndex
}
