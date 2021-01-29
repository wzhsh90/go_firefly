package company

import (
	"firefly/config"
	models "firefly/model"
	"firefly/utils"
)

type mapper struct {
	ListCount      func(name string) (int64, error)                                       `args:"name"`
	List           func(name string, pageIndex, pageSize int64) ([]models.Company, error) `args:"name,pageIndex,pageSize"`
	ExistName      func(name string) (int64, error)                                       `args:"name"`
	InsertTemplete func(arg models.Company) (int64, error)
	Update         func(arg models.Company) (int64, error)
	Del            func(id string) (int64, error)          `args:"id"`
	Get            func(id string) (models.Company, error) `args:"id"`
}

var dao mapper

func init() {
	config.RegisterMapper("resource/mybatis/CompanyMapper.xml", &dao)
}

type Service struct {
	dao *mapper
}

func (u *Service) ExistName(name string) (int64, error) {
	return dao.ExistName(name)
}
func (u *Service) InsertTemplete(arg *models.Company) (int64, error) {
	return dao.InsertTemplete(*arg)
}

func (u *Service) List(name string, pageIndex, pageSize int64) models.PageModelLay {
	var tableJsonData = models.PageModelLay{}
	var realName = utils.SqlLike(name)
	itemsCount, _ := dao.ListCount(realName)
	tableJsonData.BuildPageInfo(pageIndex, pageSize, itemsCount)
	list, _ := dao.List(realName, tableJsonData.PageIndex, tableJsonData.PageSize)

	tableJsonData.Rows = list
	return tableJsonData
}
func (u *Service) Update(arg *models.Company) (int64, error) {

	return dao.Update(*arg)
}
func (u *Service) Del(id string) (int64, error) {
	return dao.Del(id)
}
func (u *Service) Get(id string) (models.Company, error) {
	return dao.Get(id)
}
