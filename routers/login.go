package routers

import (
	"firefly/consts"
	"firefly/controller"
	"firefly/event"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	_ = event.Subscribe(consts.ROUTER_INIT_EVENT, func(r *gin.Engine) {
		ctl := controller.LoginController{}
		r.POST("/login_api", ctl.LoginApi)
		r.POST("/logoff", ctl.Logoff)
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title": "登录",
			})
		})
		r.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title": "登录",
			})
		})
		r.GET("/home", ctl.Home)

	})
}
