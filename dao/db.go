package dao

import (
	models "firefly/model"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"time"
)

var session db.Session

func InitFromConfig(dbConfig models.DataSource) {
	var settings = mysql.ConnectionURL{
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Database: dbConfig.Database,
		Host:     dbConfig.Host,
		Socket:   dbConfig.Socket,
		Options: map[string]string{
			"charset":   "utf8mb4",
			"parseTime": "true",
			"loc":       "Local",
		},
	}
	db.DefaultSettings.SetConnMaxLifetime(time.Duration(dbConfig.MaxLife) * time.Second)
	db.DefaultSettings.SetMaxIdleConns(dbConfig.MaxIdle)
	db.DefaultSettings.SetMaxOpenConns(dbConfig.MaxOpen)
	db.LC().SetLevel(db.LogLevelTrace)
	sess, err := mysql.Open(settings)
	if err != nil {
		panic(err)
	}
	session = sess

}
func Exists(table string, query map[string]interface{}) (bool, error) {
	cond := db.Cond{}
	for k, v := range query {
		cond[k] = v
	}
	return session.Collection(table).Find(cond).Exists()
}
func Save(table string, entity map[string]interface{}) (int64, error) {
	ret, err := session.SQL().InsertInto(table).Values(entity).Exec()
	if err != nil {
		return 0, err
	} else {
		return ret.RowsAffected()
	}
}
func Page(table string, cond db.Cond, listCol []string, pageIndex, pageSize uint) models.PageModelLay {
	cols := make([]interface{}, 0)
	for _, v := range listCol {
		cols = append(cols, v)
	}
	res := session.Collection(table).Find(cond).Select(cols...)
	p := res.Paginate(pageSize)
	itemsCount, _ := p.Count()
	list := make([]map[string]string, 0)
	//var list []models.Company
	//var list []interface{}
	perr := p.Page(pageIndex).All(&list)
	if perr != nil {
		println(perr.Error())
	}
	var tableJsonData = models.PageModelLay{}
	pages, _ := p.TotalPages()
	tableJsonData.Pages = pages
	tableJsonData.Page = pageIndex
	tableJsonData.PageSize = pageSize
	tableJsonData.Records = uint(itemsCount)
	tableJsonData.Rows = list
	return tableJsonData
}

func Count(table string, query map[string]interface{}) (uint64, error) {
	cond := db.Cond{}
	for k, v := range query {
		cond[k] = v
	}
	return session.Collection(table).Find(cond).Count()
}
func Del(table string, query map[string]interface{}) (int64, error) {
	ret, err := session.SQL().DeleteFrom(table).Where(query).Exec()
	if err != nil {
		return 0, err
	} else {
		return ret.RowsAffected()
	}
}
func GetCol(table string, query map[string]interface{}, listCol []string) map[string]string {

	cond := db.Cond{}
	for k, v := range query {
		cond[k] = v
	}
	cols := make([]interface{}, 0)
	for _, v := range listCol {
		cols = append(cols, v)
	}
	info := make(map[string]string)
	session.Collection(table).Find(cond).Select(cols...).One(&info)
	return info
}
func Get(table string, query map[string]interface{}) map[string]string {
	cond := db.Cond{}
	for k, v := range query {
		cond[k] = v
	}
	info := make(map[string]string)
	session.Collection(table).Find(cond).One(&info)
	return info
}
func Update(table string, query map[string]interface{}, updateItem map[string]interface{}) (int64, error) {
	cond := db.Cond{}
	for k, v := range query {
		cond[k] = v
	}
	ret, err := session.SQL().Update(table).Set(updateItem).Where(cond).Exec()
	if err != nil {
		return 0, err
	} else {
		return ret.RowsAffected()
	}
}
func CloseDb() {
	session.Close()
}
