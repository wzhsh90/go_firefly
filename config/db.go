package config

import (
	"github.com/zhuxiujia/GoMybatis"
	"io/ioutil"
	"time"
)

type DbEngine struct {
	Engine *GoMybatis.GoMybatisEngine
}

var AppDbEngine DbEngine

func initFromConfig(dbConfig DataSource) {
	myBatisEngine := GoMybatis.GoMybatisEngine{}.New()
	myBatisEngine.SetLogEnable(dbConfig.LogEnable)
	db, err := myBatisEngine.Open(dbConfig.Dialect, dbConfig.Url)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxLife) * time.Second)
	db.SetMaxIdleConns(dbConfig.MaxIdle)
	db.SetMaxOpenConns(dbConfig.MaxOpen)
	AppDbEngine = DbEngine{Engine: &myBatisEngine}
}
func RegisterMapper(mapperPath string, daoPtr interface{}) {
	xmlBytes, _ := ioutil.ReadFile(mapperPath)
	AppDbEngine.Engine.WriteMapperPtr(daoPtr, xmlBytes)
}
