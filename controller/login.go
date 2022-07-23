package controller

import (
	"firefly/middleware"
	models "firefly/model"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
}

func (c *LoginController) LoginApi(ctx *gin.Context) {
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	var rest = models.RestResult{}
	rest.Code = 1
	rest.Message = "用户信息不存在"
	if account == "" {
		rest.Message = "账号不能为空"
		ctx.JSON(200, rest)
		return
	}
	if password == "" {
		rest.Message = "密码不能为空"
		ctx.JSON(200, rest)
		return
	}
	user := models.LoginUser{
		Id:   account,
		Name: password,
	}
	rest.Code = 0
	rest.Message = "登录成功"
	rest.Result = "/home"
	middleware.SaveCurrentUser(ctx, &user)
	ctx.JSON(200, rest)
}
func (c *LoginController) Logoff(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 0
	middleware.RemoveUserSession(ctx)
	ctx.JSON(200, rest)
}
func (c *LoginController) Home(ctx *gin.Context) {
	ctx.HTML(200, "home.html", gin.H{
		"title": "首页",
	})
}
