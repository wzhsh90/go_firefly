package config

import (
	"firefly/dao"
	models "firefly/model"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var bootConfig models.BootConfig

func init() {
	envParam := flag.String("env", "dev", "go run main.go --env dev/prod")
	flag.Parse()
	var configName = ""
	if *envParam != "" {
		configName = *envParam
		if *envParam != "dev" && *envParam != "prod" {
			configName = "dev"
		}
	}
	data, cerr := ioutil.ReadFile("resource/" + configName + ".yml")
	if cerr != nil {
		log.Fatal("配置文件不存在")
	}
	//把yaml形式的字符串解析成struct类型
	_ = yaml.Unmarshal(data, &bootConfig)
	if bootConfig.Db.Dialect != "" {
		dao.InitFromConfig(bootConfig.Db)
	} else {
		log.Print("当前无数据配置")
	}
}
func GetAppConfig() models.BootConfig {

	return bootConfig
}
