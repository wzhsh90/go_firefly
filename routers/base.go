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
		ctl := controller.BaseController{}
		v1 := r.Group("/company")
		{
			v1.POST("/list", ctl.List)
			v1.POST("/add", ctl.Add)
			v1.POST("/del", ctl.Del)
			v1.POST("/update", ctl.Update)
			v1.GET("/index", func(c *gin.Context) {
				c.HTML(http.StatusOK, "company/index.html", gin.H{
					"title": "公司列表",
				})
			})

		}

	})
}
