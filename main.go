package main

import (
	"firefly/config"
	"firefly/db"
	"firefly/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	defer db.CloseDb()
	bootConfig := config.GetAppConfig()
	gin.SetMode(bootConfig.Server.Mode)
	r := gin.Default()
	//不显示gin 控制台日志时用下面方式
	//r := gin.New()
	routers.InitRouter(r, bootConfig.Server)
}
