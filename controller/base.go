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
	opeKey, ok := c.checkOpeKey(ctx.PostForm("opeKey"))
	if ok {
		return
	}
	crudInfo := models.LoadCrudFile(crudJson)
	list := crudInfo.List[opeKey]
	if list.Disable {
		ctx.String(200, "未启用查询,请检查配置文件")
		return
	}
	pageIndex := utils.ParseInt(ctx.PostForm("pageIndex"))
	pageSize := utils.ParseInt(ctx.PostForm("pageSize"))
	columMap := crudInfo.Mod.Columns
	valid := list.UnStrictParse(columMap, ctx)
	if !valid {
		ctx.String(200, "数据不合法")
		return
	}
	tableJson := dao.PageSql(crudInfo.Mod.Table.Name, list, pageIndex, pageSize)
	ctx.JSON(200, tableJson)
}

func (c *BaseController) checkOpeKey(opeKey string) (string, bool) {
	if len(opeKey) <= 0 {
		return "未启用查询,请检查配置文件", true
	}
	return opeKey, false
}
func (c *BaseController) Add(ctx *gin.Context) {
	opeKey, ok := c.checkOpeKey(ctx.PostForm("opeKey"))
	if ok {
		return
	}
	var rest = models.RestResult{}
	rest.Code = 1
	crudInfo := models.LoadCrudFile(crudJson)
	add := crudInfo.Add[opeKey]
	if add.Disable {
		ctx.String(200, "未启用新增,请检查配置文件")
		return
	}
	columMap := crudInfo.Mod.Columns

	_, dbData, validResp := add.GetFormData(columMap, ctx, true)
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
	opeKey, ok := c.checkOpeKey(ctx.PostForm("opeKey"))
	if ok {
		return
	}
	var rest = models.RestResult{}
	rest.Code = 1
	crudInfo := models.LoadCrudFile(crudJson)
	update := crudInfo.Update[opeKey]
	if update.Disable {
		ctx.String(200, "未启用修改,请检查配置文件")
		return
	}
	columMap := crudInfo.Mod.Columns
	formData, dbData, validResp := update.GetFormData(columMap, ctx, false)
	if !validResp.Valid {
		rest.Message = validResp.Msg
		ctx.JSON(200, rest)
		return
	}
	valid := update.StrictParse(columMap, ctx)
	if !valid {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
	}
	orgInfo := dao.GetColSql(crudInfo.Mod.Table.Name, update.Where, update.Select)
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
	_, serr := dao.UpdateSql(crudInfo.Mod.Table.Name, update.Where, dbData)
	if serr == nil {
		rest.Code = 0
		rest.Message = "修改成功"
	} else {
		rest.Message = "修改失败"
	}
	ctx.JSON(200, rest)
}
func (c *BaseController) Del(ctx *gin.Context) {
	opeKey, ok := c.checkOpeKey(ctx.PostForm("opeKey"))
	if ok {
		return
	}
	var rest = models.RestResult{}
	rest.Code = 1
	crudInfo := models.LoadCrudFile(crudJson)
	del := crudInfo.Del[opeKey]
	if del.Disable {
		ctx.String(200, "未启用删除,请检查配置文件")
		return
	}
	columMap := crudInfo.Mod.Columns
	valid := del.StrictParse(columMap, ctx)
	if !valid {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
	}
	if len(del.Select) != 0 {
		orgInfo := dao.GetColSql(crudInfo.Mod.Table.Name, del.Where, del.Select)
		if len(orgInfo) == 0 {
			rest.Message = "数据获取失败或已不存在"
			ctx.JSON(200, rest)
			return
		}
	}
	//真删除
	if !del.Fake {
		db, _ := dao.DelSql(crudInfo.Mod.Table.Name, del.Where)
		if db == 1 {
			rest.Code = 0
			rest.Message = "删除成功"
		} else {
			rest.Message = "删除失败"
		}

	} else {
		//假删除
		dbData := map[string]interface{}{
			"del_flag": 1,
		}
		//db, _ := dao.UpdateSql(crudInfo.Mod.Table.Name, crudInfo.Update.Where, dbData)
		db, _ := dao.UpdateSql(crudInfo.Mod.Table.Name, del.Where, dbData)
		if db == 1 {
			rest.Code = 0
			rest.Message = "删除成功"
		} else {
			rest.Message = "删除失败"
		}
	}
	ctx.JSON(200, rest)
}
