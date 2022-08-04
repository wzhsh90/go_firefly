package routers

import (
	"firefly/consts"
	"firefly/controller"
	"firefly/event"
	"firefly/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	_ = event.Subscribe(consts.ROUTER_INIT_EVENT, func(r *gin.Engine) {
		ctl := controller.BaseController{}
		v1 := r.Group("/crud")
		v1.Use(middleware.CrudAuth())
		{
			v1.POST("/handler", ctl.Handler)
		}
	})
}
