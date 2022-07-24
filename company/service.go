package company

var tableName = "sys_company_t"
var listCol = []string{"id", "com_name", "com_desc"}

type Service struct{}

//func (u *Service) BatchSave(arg []models.Company) {
//	batcher := db2.Session.SQL().InsertInto(tableName).Batch(400)
//	go func() {
//		defer batcher.Done()
//		for i := range arg {
//			batcher.Values(arg[i])
//		}
//	}()
//	batcher.Wait()
//
//}
//
//func (u *Service) List(name string, pageIndex, pageSize uint) models.PageModelLay {
//
//	cond := db.Cond{}
//	if name != "" {
//		cond["com_name"] = db.Like(utils.SqlLike(name))
//	}
//	res := dao.Session.Collection(tableName).Find(cond).Select(listCol)
//	p := res.Paginate(pageSize)
//	itemsCount, _ := p.Count()
//	//
//	var list []models.Company
//	//p.Page(pageIndex).All(&list)
//
//	var tableJsonData = models.PageModelLay{}
//	//tableJsonData.BuildPageInfo(pageIndex, pageSize, uint(itemsCount))
//	tableJsonData.BuildPageInfo(pageIndex, pageSize, 0)
//	tableJsonData.Rows = list
//	return tableJsonData
//}
//func (u *Service) List2(name string, pageIndex, pageSize uint) models.PageModelLay {
//
//	cond := db.Cond{}
//	if name != "" {
//		cond["com_name"] = db.Like(utils.SqlLike(name))
//	}
//	itemsCount, _ := db2.Session.Collection(tableName).Find(cond).Count()
//	//
//	//p := db2.Session.SQL().SelectFrom(tableName).Columns(listCol).Where(cond).Paginate(pageSize)
//	//var list []models.Company
//	//p.Page(pageIndex).All(&list)
//
//	var tableJsonData = models.PageModelLay{}
//	//tableJsonData.BuildPageInfo(pageIndex, pageSize, uint(itemsCount))
//	//tableJsonData.Rows = list
//	return tableJsonData
//}

func (u *Service) ListSql() {
	//rows, _ := db2.Session.SQL().Query(`SELECT * FROM accounts WHERE last_name = ?`, "Smith")
	//var companies []models.Company
	//iter := db2.Session.SQL().NewIterator(rows)
	//_ = iter.All(&companies)

}
