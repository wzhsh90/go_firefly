package config

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"time"
)

var DbSession db.Session

func initFromConfig(dbConfig DataSource) {
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
	DbSession = sess

}
func CloseDb() {
	DbSession.Close()
}
