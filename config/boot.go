package config

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type BootConfig struct {
	Server Server
	Db     DataSource
}
type Server struct {
	Mode        string
	Port        string
	SessionName string
	SessionKey  string
	TplPath     string
}
type DataSource struct {
	Dialect   string
	User      string
	Password  string
	Database  string
	Host      string
	Socket    string
	MaxLife   int
	MaxIdle   int
	MaxOpen   int
	LogEnable bool
}

var bootConfig BootConfig

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
		initFromConfig(bootConfig.Db)
	} else {
		log.Print("当前无数据配置")
	}
}
func GetAppConfig() BootConfig {

	return bootConfig
}
