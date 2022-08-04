package routers

import (
	"firefly/consts"
	"firefly/event"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	_ = event.Subscribe(consts.ROUTER_INIT_EVENT, func(r *gin.Engine) {
		v1 := r.Group("/company")
		{
			v1.GET("/index", func(c *gin.Context) {
				c.HTML(http.StatusOK, "company/index.html", gin.H{
					"title": "公司列表",
				})
			})

		}

	})
}
