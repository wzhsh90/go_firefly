package controller

import (
	"firefly/dao"
	models "firefly/model"
	"firefly/utils"
	"github.com/gin-gonic/gin"
)

var modJson = "resource/mod/company/company.mod.json"
var listJson = "resource/mod/company/company.list.json"
var delJson = "resource/mod/company/company.del.json"
var addJson = "resource/mod/company/company.add.json"
var updateJson = "resource/mod/company/company.update.json"

type BaseController struct {
}

func (c *BaseController) List(ctx *gin.Context) {
	formQuery, modInfo := models.FormQueryLoadFile(listJson)
	pageIndex := utils.ParseInt(ctx.PostForm("pageIndex"))
	pageSize := utils.ParseInt(ctx.PostForm("pageSize"))
	columMap := modInfo.Columns
	valid := formQuery.UnStrictParse(columMap, ctx)
	if !valid {
		ctx.String(200, "数据不合法")
		return
	}
	tableJson := dao.PageSql(modInfo.Table.Name, formQuery.Where, formQuery.Select, pageIndex, pageSize)
	ctx.JSON(200, tableJson)
}
func (c *BaseController) Add(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	formAdd, modInfo := models.FormAddLoadFile(addJson)
	columMap := modInfo.Columns
	_, dbData, validResp := formAdd.GetFormData(columMap, ctx, true)
	if !validResp.Valid {
		rest.Message = validResp.Msg
		ctx.JSON(200, rest)
		return
	}
	//todo 处理File 上传问题
	//判断当前数据是否存在
	if len(formAdd.Exits.Columns) >= 1 {
		existMap := make(map[string]interface{})
		for _, v := range formAdd.Exits.Columns {
			existMap[v] = dbData[v]
		}
		existFlag, _ := dao.Exists(modInfo.Table.Name, existMap)
		if existFlag {
			rest.Message = formAdd.Exits.Tip
			ctx.JSON(200, rest)
			return
		}
	}
	_, serr := dao.Save(modInfo.Table.Name, dbData)
	if serr == nil {
		rest.Code = 0
		rest.Message = "添加成功"
	} else {
		rest.Message = "添加失败"
	}
	ctx.JSON(200, rest)
}
func (c *BaseController) Update(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	formUpdate, modInfo := models.FormUpdateLoadFile(updateJson)
	columMap := modInfo.Columns
	formData, dbData, validResp := formUpdate.GetFormData(columMap, ctx, false)
	if !validResp.Valid {
		rest.Message = validResp.Msg
		ctx.JSON(200, rest)
		return
	}
	valid := formUpdate.StrictParse(columMap, ctx)
	if !valid {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
	}
	orgInfo := dao.GetColSql(modInfo.Table.Name, formUpdate.Where, formUpdate.Select)
	if len(orgInfo) == 0 {
		rest.Message = "获取历史数据失败或已不存在"
		ctx.JSON(200, rest)
		return
	}
	//判断当前数据是否存在
	if len(formUpdate.Exits.Columns) >= 1 {
		existMap := make(map[string]interface{})
		checkFlag := false
		for _, v := range formUpdate.Exits.Columns {
			existMap[v] = formData[v]
			if orgInfo[v] != formData[v] {
				checkFlag = true
				break
			}
			//interface conversion: interface {} is []uint8, not string
			if orgInfo[v] != formData[v] {
				checkFlag = true
				break
			}
		}
		if checkFlag {
			existFlag, _ := dao.Exists(modInfo.Table.Name, existMap)
			if existFlag {
				rest.Message = formUpdate.Exits.Tip
				ctx.JSON(200, rest)
				return
			}
		}
	}
	_, serr := dao.UpdateSql(modInfo.Table.Name, formUpdate.Where, dbData)
	if serr == nil {
		rest.Code = 0
		rest.Message = "修改成功"
	} else {
		rest.Message = "修改失败"
	}
	ctx.JSON(200, rest)
}
func (c *BaseController) Del(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	formQuery, modInfo := models.FormQueryLoadFile(delJson)
	columMap := modInfo.Columns
	valid := formQuery.StrictParse(columMap, ctx)
	if !valid {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
	}
	db, _ := dao.DelSql(modInfo.Table.Name, formQuery.Where)
	if db == 1 {
		rest.Code = 0
		rest.Message = "删除成功"
	} else {
		rest.Message = "删除失败"
	}
	ctx.JSON(200, rest)
}
