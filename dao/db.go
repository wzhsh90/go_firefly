package dao

import (
	models "firefly/model"
	"firefly/utils"
	"github.com/antonmedv/expr"
	"github.com/huandu/go-sqlbuilder"
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
func SelectBuilder() *sqlbuilder.SelectBuilder {
	sb := sqlbuilder.NewSelectBuilder()
	return sb
}
func PageSql(table string, cond []models.FormQueryOp, listCol []string, pageIndex, pageSize int) models.PageModelLay {
	var itemsCount int64
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("count(*)").From(table)
	if len(cond) >= 1 {
		for _, v := range cond {
			expFlag := false
			if v.On != "" {
				env := map[string]interface{}{
					v.Name: v.Val,
				}
				program, err := expr.Compile(v.On, expr.Env(env), expr.AsBool())
				if err == nil {
					output, rerr := expr.Run(program, env)
					if rerr == nil {
						expFlag = output.(bool)
					}
				}
			} else {
				expFlag = true
			}
			if !expFlag {
				continue
			}
			if v.Op == "eq" {
				sb.Where(sb.Equal(v.Name, v.Val))
			} else if v.Op == "like" {
				sb.Where(sb.Like(v.Name, utils.SqlLike(v.Val.(string))))
			} else if v.Op == "llike" {
				sb.Where(sb.Like(v.Name, utils.LeftLike(v.Val.(string))))
			} else if v.Op == "rlike" {
				sb.Where(sb.Like(v.Name, utils.RightLike(v.Val.(string))))
			} else if v.Op == "ge" {
				sb.Where(sb.GE(v.Name, v.Val))
			} else if v.Op == "gt" {
				sb.Where(sb.G(v.Name, v.Val))
			} else if v.Op == "le" {
				sb.Where(sb.LE(v.Name, v.Val))
			} else if v.Op == "lt" {
				sb.Where(sb.L(v.Name, v.Val))
			}
		}
	}
	countSql, countArgs := sb.Build()
	//println(countSql)
	session.Raw(countSql, countArgs...).Scan(&itemsCount)
	var tableJsonData = models.PageModelLay{}
	tableJsonData.BuildPageInfo(pageIndex, pageSize, int(itemsCount))
	list := make([]map[string]interface{}, 0)
	sb.Select(listCol...).From(table)
	listSql, _ := sb.Offset(tableJsonData.PageIndex).Limit(tableJsonData.PageSize).Build()
	//println(listSql)
	session.Raw(listSql, countArgs...).Find(&list)
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
