package company

import (
	"firefly/config"
	models "firefly/model"
	"firefly/utils"
	"github.com/upper/db/v4"
)

var tableName = "sys_company_t"
var listCol = []string{"id", "com_name", "com_desc"}

type Service struct{}

func (u *Service) ExistName(name string) (bool, error) {

	return config.DbSession.Collection(tableName).Find("com_name=?", name).Exists()
}
func (u *Service) Save(arg models.Company) error {
	ret, err := config.DbSession.SQL().InsertInto(tableName).Values(arg).Exec()
	if err != nil {
		return err
	} else {
		_, rerr := ret.RowsAffected()
		return rerr
	}

	//return config.DbSession.Save(&arg)
}
func (u *Service) BatchSave(arg []models.Company) {
	batcher := config.DbSession.SQL().InsertInto(tableName).Batch(400)
	go func() {
		defer batcher.Done()
		for i := range arg {
			batcher.Values(arg[i])
		}
	}()
	batcher.Wait()

}

func (u *Service) List(name string, pageIndex, pageSize uint) models.PageModelLay {

	cond := db.Cond{}
	if name != "" {
		cond["com_name"] = db.Like(utils.SqlLike(name))
	}
	res := config.DbSession.Collection(tableName).Find(cond).Select(listCol)
	p := res.Paginate(pageSize)
	itemsCount, _ := p.Count()

	var list []models.Company
	p.Page(pageIndex).All(&list)

	var tableJsonData = models.PageModelLay{}
	tableJsonData.BuildPageInfo(pageIndex, pageSize, uint(itemsCount))
	tableJsonData.Rows = list
	return tableJsonData
}
func (u *Service) List2(name string, pageIndex, pageSize uint) models.PageModelLay {

	cond := db.Cond{}
	if name != "" {
		cond["com_name"] = db.Like(utils.SqlLike(name))
	}
	itemsCount, _ := config.DbSession.Collection(tableName).Find(cond).Count()

	p := config.DbSession.SQL().SelectFrom(tableName).Columns(listCol).Where(cond).Paginate(pageSize)
	var list []models.Company
	p.Page(pageIndex).All(&list)

	var tableJsonData = models.PageModelLay{}
	tableJsonData.BuildPageInfo(pageIndex, pageSize, uint(itemsCount))
	tableJsonData.Rows = list
	return tableJsonData
}

func (u *Service) ListSql() {
	rows, _ := config.DbSession.SQL().Query(`SELECT * FROM accounts WHERE last_name = ?`, "Smith")
	var companies []models.Company
	iter := config.DbSession.SQL().NewIterator(rows)
	_ = iter.All(&companies)

}
func (u *Service) Update(arg models.Company) error {

	return config.DbSession.Collection(tableName).Find("id=?", arg.Id).Update(&arg)
}
func (u *Service) Del(id string) (int64, error) {

	//1、config.DbSession.Delete(&models.Company{Id: arg.Id})
	//2、config.DbSession.Collection(tableName).Find("id=?", id).Delete()

	ret, err := config.DbSession.SQL().DeleteFrom(tableName).Where("id=?", id).Exec()
	if err != nil {
		return 0, err
	} else {
		return ret.RowsAffected()
	}
}
func (u *Service) Get(id string) models.Company {
	var entity models.Company
	config.DbSession.Collection(tableName).Find("id", id).One(&entity)
	//config.DbSession.Get(&entity, db.Cond{"id": id})
	return entity
}
