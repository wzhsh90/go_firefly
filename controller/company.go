package controller

import (
	"firefly/company"
	models "firefly/model"
	"firefly/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CompanyController struct {
}

func (c *CompanyController) List(ctx *gin.Context) {
	name := ctx.PostForm("name")
	pageIndex, _ := strconv.ParseInt(ctx.PostForm("pageIndex"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.PostForm("pageSize"), 10, 64)
	svr := company.Service{}
	tableJson := svr.List(name, pageIndex, pageSize)
	ctx.JSON(200, tableJson)
}
func (c *CompanyController) Add(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	comeName := ctx.PostForm("com_name")
	if comeName == "" {
		rest.Message = "公司名称不能为空"
		ctx.JSON(200, rest)
		return
	}
	comDesc := ctx.PostForm("com_desc")
	svr := company.Service{}
	existFlag, _ := svr.ExistName(comeName)
	if existFlag >= 1 {
		rest.Message = "当前公司名称已经存在"
		ctx.JSON(200, rest)
		return
	}
	dbm := models.Company{
		ComDesc: comDesc,
		ComName: comeName,
		Id:      utils.NewId(),
	}
	db, _ := svr.InsertTemplete(&dbm)
	if db == 1 {
		rest.Code = 0
		rest.Message = "添加成功"
	} else {
		rest.Message = "添加失败"
	}
	ctx.JSON(200, rest)
}
func (c *CompanyController) Update(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	comeName := ctx.PostForm("com_name")
	if comeName == "" {
		rest.Message = "公司名称不能为空"
		ctx.JSON(200, rest)
		return
	}
	comDesc := ctx.PostForm("com_desc")
	svr := company.Service{}
	org, oerr := svr.Get(ctx.PostForm("id"))
	if oerr != nil {
		rest.Message = "当前数据不存在"
		ctx.JSON(200, rest)
		return
	}
	if org.ComName != comeName {
		existFlag, _ := svr.ExistName(comeName)
		if existFlag >= 1 {
			rest.Message = "当前公司名称已经存在"
			ctx.JSON(200, rest)
			return
		}
	}
	dbm := models.Company{
		ComDesc: comDesc,
		ComName: comeName,
		Id:      ctx.PostForm("id"),
	}
	if _, err := svr.Update(&dbm); err == nil {
		rest.Code = 0
		rest.Message = "修改成功"
	} else {
		rest.Message = "修改失败"
	}
	ctx.JSON(200, rest)
}
func (c *CompanyController) Del(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	id := ctx.PostForm("id")
	if id == "" {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
		return
	}
	svr := company.Service{}
	db, _ := svr.Del(id)
	if db == 1 {
		rest.Code = 0
		rest.Message = "删除成功"
	} else {
		rest.Message = "删除失败"
	}
	ctx.JSON(200, rest)
}
