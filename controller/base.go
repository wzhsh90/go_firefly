package controller

import (
	"firefly/dao"
	models "firefly/model"
	"firefly/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

//var crudJson = "resource/mod/company/company.crud.json"
var processMap = map[string]func(cfg models.CrudInfo, ctx *gin.Context, node string){}

func registerProcess(name string, handler func(cfg models.CrudInfo, ctx *gin.Context, node string)) {
	name = strings.ToLower(name)
	processMap[name] = handler
}
func init() {
	registerProcess("page", pageProcess)
	registerProcess("list", listProcess)
	registerProcess("get", getProcess)
	registerProcess("add", addProcess)
	registerProcess("update", updateProcess)
	registerProcess("del", delProcess)
}

type BaseController struct {
}

func (c *BaseController) Handler(ctx *gin.Context) {
	code := ctx.PostForm("code")
	engine := models.NewCrudEngine(code)
	cfg, cfgOk := models.LoadCrudFile(engine.Name)
	if !cfgOk {
		ctx.String(401, "参数非法，权限不存在")
		return
	}
	//处理器
	if processFn, ok := processMap[engine.Process]; ok {
		processFn(cfg, ctx, engine.Node)
	} else {
		ctx.String(401, "参数非法，权限不存在")
		return
	}

}
func pageProcess(crudInfo models.CrudInfo, ctx *gin.Context, node string) {
	page, ok := crudInfo.Page[node]
	if !ok || page.Disable {
		ctx.String(200, "未启用查询,请检查配置文件")
		return
	}
	pageIndex := utils.ParseInt(ctx.PostForm("pageIndex"))
	pageSize := utils.ParseInt(ctx.PostForm("pageSize"))
	columMap := crudInfo.Mod.Columns
	valid := page.UnStrictParse(columMap, ctx)
	if !valid {
		ctx.String(200, "数据不合法")
		return
	}
	tableJson := dao.PageSql(page.From, page, pageIndex, pageSize)
	ctx.JSON(200, tableJson)
}
func getProcess(crudInfo models.CrudInfo, ctx *gin.Context, node string) {
	get, ok := crudInfo.Get[node]
	if !ok || get.Disable {
		ctx.String(200, "未启用查询,请检查配置文件")
		return
	}
	columMap := crudInfo.Mod.Columns
	valid := get.UnStrictParse(columMap, ctx)
	if !valid {
		ctx.String(200, "数据不合法")
		return
	}
	tableJson := dao.GetColSql(get.From, get.Where, get.Select)
	ctx.JSON(200, tableJson)
}
func listProcess(crudInfo models.CrudInfo, ctx *gin.Context, node string) {
	list, ok := crudInfo.List[node]
	if !ok || list.Disable {
		ctx.String(200, "未启用查询,请检查配置文件")
		return
	}
	columMap := crudInfo.Mod.Columns
	valid := list.UnStrictParse(columMap, ctx)
	if !valid {
		ctx.String(200, "数据不合法")
		return
	}
	tableJson := dao.ListColSql(list.From, list.Where, list.Select)
	ctx.JSON(200, tableJson)
}
func addProcess(crudInfo models.CrudInfo, ctx *gin.Context, node string) {
	var rest = models.RestResult{}
	rest.Code = 1
	add, ok := crudInfo.Add[node]
	if !ok || add.Disable {
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
func updateProcess(crudInfo models.CrudInfo, ctx *gin.Context, node string) {
	var rest = models.RestResult{}
	rest.Code = 1
	update, ok := crudInfo.Update[node]
	if !ok || update.Disable {
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
func delProcess(crudInfo models.CrudInfo, ctx *gin.Context, node string) {
	var rest = models.RestResult{}
	rest.Code = 1
	del, ok := crudInfo.Del[node]
	if !ok || del.Disable {
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
