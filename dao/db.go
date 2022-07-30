package dao

import (
	models "firefly/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var session *gorm.DB

type EmptyStruct struct {
}

func InitFromConfig(dbConfig models.DataSource) {
	db, err := gorm.Open(mysql.Open(dbConfig.Url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	session = db

}
func Exists(table string, query map[string]interface{}) (bool, error) {
	var totalCnt int64
	tx := session.Table(table).Where(query).Count(&totalCnt)
	return totalCnt >= 1, tx.Error
}
func Save(table string, entity map[string]interface{}) (int64, error) {
	tx := session.Table(table).Create(entity)
	return tx.RowsAffected, tx.Error
}
func Page(table string, cond map[string]interface{}, listCol []string, pageIndex, pageSize int) models.PageModelLay {
	var itemsCount int64
	_ = session.Table(table).Where(cond).Count(&itemsCount)
	list := make([]map[string]interface{}, 0)
	var tableJsonData = models.PageModelLay{}
	tableJsonData.BuildPageInfo(pageIndex, pageSize, int(itemsCount))
	session.Table(table).Where(cond).Select(listCol).Offset(tableJsonData.PageIndex).Limit(tableJsonData.PageSize).Find(&list)
	tableJsonData.Rows = list
	return tableJsonData
}
func PageStruct(table string, cond map[string]interface{}, listCol []string, listPtr interface{}, pageIndex, pageSize int) models.PageModelLay {
	var itemsCount int64
	_ = session.Table(table).Where(cond).Count(&itemsCount)
	var tableJsonData = models.PageModelLay{}
	tableJsonData.BuildPageInfo(pageIndex, pageSize, int(itemsCount))
	session.Table(table).Where(cond).Select(listCol).Offset(tableJsonData.PageIndex).Limit(tableJsonData.PageSize).Find(listPtr)
	tableJsonData.Rows = listPtr
	return tableJsonData
}

func Count(table string, query map[string]interface{}) (int64, error) {
	var totalCnt int64
	tx := session.Table(table).Where(query).Count(&totalCnt)
	return totalCnt, tx.Error
}
func Del(table string, query map[string]interface{}) (int64, error) {
	info := EmptyStruct{}
	tx := session.Table(table).Where(query).Delete(&info)
	return tx.RowsAffected, tx.Error
}
func GetCol(table string, query map[string]interface{}, listCol []string) map[string]interface{} {
	info := make(map[string]interface{})
	session.Table(table).Where(query).Select(listCol).Limit(1).Find(&info)
	//for k, v := range info {
	//	switch v.(type) {
	//	case []uint8:
	//		arr := v.([]uint8)
	//		info[k] = string(arr)
	//	case nil:
	//		info[k] = ""
	//	}
	//}
	return info
}

func GetColStruct(table string, query map[string]interface{}, listCol []string, entityPtr interface{}) {
	session.Table(table).Where(query).Select(listCol).Limit(1).Find(entityPtr)
}
func Get(table string, query map[string]interface{}) map[string]interface{} {
	info := make(map[string]interface{})
	session.Table(table).Where(query).Limit(1).Find(&info)
	//for k, v := range info {
	//	switch v.(type) {
	//	case []uint8:
	//		arr := v.([]uint8)
	//		info[k] = string(arr)
	//	case nil:
	//		info[k] = ""
	//	}
	//}
	return info
}

func GetStruct(table string, query map[string]interface{}, entityPtr interface{}) {
	session.Table(table).Where(query).Limit(1).Find(entityPtr)
}
func Update(table string, query map[string]interface{}, updateItem map[string]interface{}) (int64, error) {
	tx := session.Table(table).Where(query).UpdateColumns(updateItem)
	return tx.RowsAffected, tx.Error
}
func CloseDb() {
}
