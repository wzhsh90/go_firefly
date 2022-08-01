package controller

import (
	"firefly/dao"
	models "firefly/model"
	"firefly/utils"
	"github.com/gin-gonic/gin"
)

var crudJson = "resource/mod/company/company.crud.json"

type BaseController struct {
}

func (c *BaseController) List(ctx *gin.Context) {
	crudInfo := models.LoadCrudFile(crudJson)
	pageIndex := utils.ParseInt(ctx.PostForm("pageIndex"))
	pageSize := utils.ParseInt(ctx.PostForm("pageSize"))
	columMap := crudInfo.Mod.Columns
	valid := crudInfo.List.UnStrictParse(columMap, ctx)
	if !valid {
		ctx.String(200, "数据不合法")
		return
	}
	tableJson := dao.PageSql(crudInfo.Mod.Table.Name, crudInfo.List.Where, crudInfo.List.Select, crudInfo.List.Order, pageIndex, pageSize)
	ctx.JSON(200, tableJson)
}
func (c *BaseController) Add(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	crudInfo := models.LoadCrudFile(crudJson)
	columMap := crudInfo.Mod.Columns
	_, dbData, validResp := crudInfo.Add.GetFormData(columMap, ctx, true)
	if !validResp.Valid {
		rest.Message = validResp.Msg
		ctx.JSON(200, rest)
		return
	}
	//todo 处理File 上传问题
	//判断当前数据是否存在
	if len(crudInfo.Mod.Unique.Columns) >= 1 {
		existMap := make(map[string]interface{})
		for _, v := range crudInfo.Mod.Unique.Columns {
			existMap[v] = dbData[v]
		}
		existFlag, _ := dao.Exists(crudInfo.Mod.Table.Name, existMap)
		if existFlag {
			rest.Message = crudInfo.Mod.Unique.Tip
			ctx.JSON(200, rest)
			return
		}
	}
	_, serr := dao.Save(crudInfo.Mod.Table.Name, dbData)
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
	crudInfo := models.LoadCrudFile(crudJson)
	columMap := crudInfo.Mod.Columns
	formData, dbData, validResp := crudInfo.Update.GetFormData(columMap, ctx, false)
	if !validResp.Valid {
		rest.Message = validResp.Msg
		ctx.JSON(200, rest)
		return
	}
	valid := crudInfo.Update.StrictParse(columMap, ctx)
	if !valid {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
	}
	orgInfo := dao.GetColSql(crudInfo.Mod.Table.Name, crudInfo.Update.Where, crudInfo.Update.Select)
	if len(orgInfo) == 0 {
		rest.Message = "数据获取失败或已不存在"
		ctx.JSON(200, rest)
		return
	}
	//判断当前数据是否存在
	if len(crudInfo.Mod.Unique.Columns) >= 1 {
		existMap := make(map[string]interface{})
		checkFlag := false
		for _, v := range crudInfo.Mod.Unique.Columns {
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
			existFlag, _ := dao.Exists(crudInfo.Mod.Table.Name, existMap)
			if existFlag {
				rest.Message = crudInfo.Mod.Unique.Tip
				ctx.JSON(200, rest)
				return
			}
		}
	}
	_, serr := dao.UpdateSql(crudInfo.Mod.Table.Name, crudInfo.Update.Where, dbData)
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
	crudInfo := models.LoadCrudFile(crudJson)
	columMap := crudInfo.Mod.Columns
	valid := crudInfo.Del.StrictParse(columMap, ctx)
	if !valid {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
	}
	if len(crudInfo.Del.Select) != 0 {
		orgInfo := dao.GetColSql(crudInfo.Mod.Table.Name, crudInfo.Del.Where, crudInfo.Del.Select)
		if len(orgInfo) == 0 {
			rest.Message = "数据获取失败或已不存在"
			ctx.JSON(200, rest)
			return
		}
	}
	if crudInfo.Del.Physic {
		db, _ := dao.DelSql(crudInfo.Mod.Table.Name, crudInfo.Del.Where)
		if db == 1 {
			rest.Code = 0
			rest.Message = "删除成功"
		} else {
			rest.Message = "删除失败"
		}

	} else {
		dbData := map[string]interface{}{
			"del_flag": 1,
		}
		db, _ := dao.UpdateSql(crudInfo.Mod.Table.Name, crudInfo.Del.Where, dbData)
		if db == 1 {
			rest.Code = 0
			rest.Message = "删除成功"
		} else {
			rest.Message = "删除失败"
		}
	}
	ctx.JSON(200, rest)
}
