package controller

import (
	"firefly/company"
	"firefly/db"
	models "firefly/model"
	"firefly/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"strconv"
)

var modJson = "resource/mod/company/company.mod.json"
var addJson = "resource/mod/company/company.add.json"

var validate = validator.New()

type BaseController struct {
}

func (c *BaseController) List(ctx *gin.Context) {
	name := ctx.PostForm("name")
	pageIndex := utils.ParseUnInt(ctx.PostForm("pageIndex"))
	pageSize := utils.ParseUnInt(ctx.PostForm("pageSize"))
	svr := company.Service{}
	tableJson := svr.List(name, pageIndex, pageSize)
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
		existFlag, _ := db.Exists(entity.Table.Name, existMap)
		if existFlag {
			rest.Message = formAdd.Exits.Tip
			ctx.JSON(200, rest)
			return
		}
	}
	_, serr := db.Save(entity.Table.Name, dbData)
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
	comeName := ctx.PostForm("com_name")
	if comeName == "" {
		rest.Message = "公司名称不能为空"
		ctx.JSON(200, rest)
		return
	}
	comDesc := ctx.PostForm("com_desc")
	svr := company.Service{}
	org := svr.Get(ctx.PostForm("id"))
	if org.ComName != comeName {
		existFlag, _ := db.Exists("", map[string]interface{}{
			"com_name": comeName,
		})
		if existFlag {
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
	if err := svr.Update(dbm); err == nil {
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
	db, _ := db.Del("", map[string]interface{}{
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
