package db

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
func Del(table string, query map[string]interface{}) (int64, error) {
	ret, err := session.SQL().DeleteFrom(table).Where(query).Exec()
	if err != nil {
		return 0, err
	} else {
		return ret.RowsAffected()
	}
}
func Update(table string, query map[string]interface{}, updateItem map[string]interface{}) (int64, error) {
	cond := db.Cond{}
	for k, v := range query {
		cond[k] = v
	}
	ret, err := session.SQL().Update(table).Set(updateItem).Where(query).Exec()
	if err != nil {
		return 0, err
	} else {
		return ret.RowsAffected()
	}
}
func CloseDb() {
	session.Close()
}
