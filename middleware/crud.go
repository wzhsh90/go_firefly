package middleware

import (
	models "firefly/model"
	"github.com/gin-gonic/gin"
	"net/http"

	"strings"
)

func invalid(c *gin.Context, msg string) {
	rest := models.RestResult{}
	rest.Code = 100
	rest.Message = msg
	c.JSON(http.StatusUnauthorized, rest)
	c.Abort()
}
func CrudAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.PostForm("code")
		if code == "" {
			invalid(ctx, "操作非法")
			return
		}
		codeList := strings.Split(code, ".")
		if len(codeList) != 3 {
			invalid(ctx, "操作非法")
			return
		}
		ctx.Next()
	}
}
