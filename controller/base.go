package controller

import (
	"firefly/dao"
	models "firefly/model"
	"firefly/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/upper/db/v4"

	"strconv"
)

var modJson = "resource/mod/company/company.mod.json"
var addJson = "resource/mod/company/company.add.json"
var updateJson = "resource/mod/company/company.update.json"
var listCol = []string{"id", "com_name", "com_desc"}
var validate = validator.New()

type BaseController struct {
}

func (c *BaseController) List(ctx *gin.Context) {
	name := ctx.PostForm("name")
	pageIndex := utils.ParseUnInt(ctx.PostForm("pageIndex"))
	pageSize := utils.ParseUnInt(ctx.PostForm("pageSize"))
	entity := models.ModInfoGenFile(modJson)
	cond := db.Cond{}
	if name != "" {
		cond["com_name"] = db.Like(utils.SqlLike(name))
	}
	tableJson := dao.Page(entity.Table.Name, cond, listCol, pageIndex, pageSize)
	ctx.JSON(200, tableJson)
}
func (c *BaseController) Add(ctx *gin.Context) {
	var rest = models.RestResult{}
	rest.Code = 1
	entity := models.ModInfoGenFile(modJson)
	formAdd := models.FormAddModGenFile(addJson)
	columMap := entity.Columns
	dbData := make(map[string]interface{})
	formData := make(map[string]interface{})

	//生成form data
	for k, formCol := range formAdd.Columns {
		val := ctx.PostForm(k)
		dbCol, dbOk := columMap[k]
		if !dbOk {
			rest.Message = "数据不合法"
			ctx.JSON(200, rest)
			return
		}
		if dbCol.LangType == "int" {
			pval, per := strconv.ParseInt(val, 10, 64)
			if per != nil {
				rest.Message = formCol.ZhName + "数据不合法"
				ctx.JSON(200, rest)
				return
			}
			formData[k] = pval
		} else if dbCol.LangType == "float" {
			pval, per := strconv.ParseFloat(val, 64)
			if per != nil {
				rest.Message = formCol.ZhName + "数据不合法"
				ctx.JSON(200, rest)
				return
			}
			formData[k] = pval
		} else {
			formData[k] = val
		}
	}
	//todo 处理File 上传问题
	//检验form 上传是否合法
	rules := formAdd.GetRules()
	if len(rules) != 0 {
		errs := validate.ValidateMap(formData, rules)
		if len(errs) > 0 {
			rest.Message = "数据不合法"
			ctx.JSON(200, rest)
			return
		}
	}
	//根据mod 生成dbdata
	for k, dbCol := range columMap {
		formVal, formOk := formData[k]
		if formOk {
			dbData[k] = formVal
		} else {
			if dbCol.LangType == "string" {
				dbData[k] = ""
				if dbCol.Key {
					dbData[k] = utils.NewId()
				}
			} else if dbCol.LangType == "int" || dbCol.LangType == "float" {
				dbData[k] = 0
			}
		}
	}
	//判断当前数据是否存在
	if len(formAdd.Exits.Columns) >= 1 {
		existMap := make(map[string]interface{})
		for _, v := range formAdd.Exits.Columns {
			existMap[v] = dbData[v]
		}
		existFlag, _ := dao.Exists(entity.Table.Name, existMap)
		if existFlag {
			rest.Message = formAdd.Exits.Tip
			ctx.JSON(200, rest)
			return
		}
	}
	_, serr := dao.Save(entity.Table.Name, dbData)
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
	entity := models.ModInfoGenFile(modJson)
	formAdd := models.FormAddModGenFile(updateJson)
	columMap := entity.Columns
	formData := make(map[string]interface{})

	//生成form data
	for k, formCol := range formAdd.Columns {
		val := ctx.PostForm(k)
		dbCol, dbOk := columMap[k]
		if !dbOk {
			rest.Message = "数据不合法"
			ctx.JSON(200, rest)
			return
		}
		if dbCol.LangType == "int" {
			pval, per := strconv.ParseInt(val, 10, 64)
			if per != nil {
				rest.Message = formCol.ZhName + "数据不合法"
				ctx.JSON(200, rest)
				return
			}
			formData[k] = pval
		} else if dbCol.LangType == "float" {
			pval, per := strconv.ParseFloat(val, 64)
			if per != nil {
				rest.Message = formCol.ZhName + "数据不合法"
				ctx.JSON(200, rest)
				return
			}
			formData[k] = pval
		} else {
			formData[k] = val
		}
	}
	//todo 处理File 上传问题
	//检验form 上传是否合法
	rules := formAdd.GetRules()
	if len(rules) != 0 {
		errs := validate.ValidateMap(formData, rules)
		if len(errs) > 0 {
			rest.Message = "数据不合法"
			ctx.JSON(200, rest)
			return
		}
	}
	keys := entity.GetKeyCols()
	queryMap := make(map[string]interface{})
	for k, _ := range keys {
		val, ok := formData[k]
		if !ok {
			rest.Message = "数据不合法"
			ctx.JSON(200, rest)
			return
		}
		queryMap[k] = val
		//删除修改主键数据,不允许修改主键
		delete(formData, k)
	}
	orgInfo := dao.Get(entity.Table.Name, queryMap)
	if len(orgInfo) == 0 {
		rest.Message = "获取历史数据失败或已不存在"
		ctx.JSON(200, rest)
		return
	}

	//判断当前数据是否存在
	if len(formAdd.Exits.Columns) >= 1 {
		existMap := make(map[string]interface{})
		checkFlag := false
		for _, v := range formAdd.Exits.Columns {
			existMap[v] = formData[v]
			if orgInfo[v] != formData[v] {
				checkFlag = true
				break
			}
		}
		if checkFlag {
			existFlag, _ := dao.Exists(entity.Table.Name, existMap)
			if existFlag {
				rest.Message = formAdd.Exits.Tip
				ctx.JSON(200, rest)
				return
			}
		}
	}
	_, serr := dao.Update(entity.Table.Name, queryMap, formData)
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
	id := ctx.PostForm("id")
	if id == "" {
		rest.Message = "数据不合法"
		ctx.JSON(200, rest)
		return
	}
	entity := models.ModInfoGenFile(modJson)
	db, _ := dao.Del(entity.Table.Name, map[string]interface{}{
		"id": id,
	})
	if db == 1 {
		rest.Code = 0
		rest.Message = "删除成功"
	} else {
		rest.Message = "删除失败"
	}
	ctx.JSON(200, rest)
}
