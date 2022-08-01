package models

import (
	"firefly/utils"
	"github.com/antonmedv/expr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

var validate = validator.New()

type ModInfo struct {
	Name    string                `json:"name"`
	Table   TableInfo             `json:"table"`
	Columns map[string]ColumnInfo `json:"columns"`
	Unique  FormUnique            `json:"unique"`
}

func (c *ModInfo) GetKeyCols() map[string]ColumnInfo {
	keys := make(map[string]ColumnInfo)
	for k, v := range c.Columns {
		if v.Key {
			keys[k] = v
		}
	}
	return keys
}

type TableInfo struct {
	Name   string `json:"name"`
	ZhName string `json:"zh_name"`
}
type ColumnInfo struct {
	ZhName   string      `json:"zh_name"`
	DbType   string      `json:"db_type"`
	LangType string      `json:"lang_type"`
	Default  interface{} `json:"default"`
	Key      bool        `json:"key"`
	Index    string      `json:"index"`
}
type FormInfo struct {
	Columns map[string]FormField `json:"columns"`
}

type FormField struct {
	ZhName string `json:"zh_name"`
	Dom    string `json:"dom"`
	Rule   string `json:"rule"`
}

func (c *FormInfo) GetRules() map[string]interface{} {
	rules := make(map[string]interface{})
	for k, v := range c.Columns {
		if v.Rule != "" {
			rules[k] = v.Rule
		}
	}
	return rules
}

type ValidResp struct {
	Valid bool
	Msg   string
}

func (c *FormInfo) GetFormData(columMap map[string]ColumnInfo, ctx *gin.Context, create bool) (map[string]interface{}, map[string]interface{}, ValidResp) {
	formData := make(map[string]interface{})
	dbData := make(map[string]interface{})
	//生成form data
	for k, formCol := range c.Columns {
		val := ctx.PostForm(k)
		dbCol, dbOk := columMap[k]
		if !dbOk {
			return nil, nil, ValidResp{
				Valid: false,
				Msg:   "数据不合法,存在未定义字段" + k,
			}
		}
		if dbCol.LangType == "int" {
			pval, per := strconv.ParseInt(val, 10, 64)
			if per != nil {
				return nil, nil, ValidResp{
					Valid: false,
					Msg:   formCol.ZhName + "数据不合法",
				}
			}
			formData[k] = pval
		} else if dbCol.LangType == "float" {
			pval, per := strconv.ParseFloat(val, 64)
			if per != nil {
				return nil, nil, ValidResp{
					Valid: false,
					Msg:   formCol.ZhName + "数据不合法",
				}
			}
			formData[k] = pval
		} else {
			formData[k] = val
		}
	}
	//检验form 上传是否合法
	rules := c.GetRules()
	if len(rules) != 0 {
		errs := validate.ValidateMap(formData, rules)
		if len(errs) > 0 {
			return nil, nil, ValidResp{
				Valid: false,
				Msg:   "数据校验不通过",
			}
		}
	}
	if create {
		for k, dbCol := range columMap {
			formVal, formOk := formData[k]
			if formOk {
				dbData[k] = formVal
			} else {
				if dbCol.Default != nil {
					//处理json string 里面是int 但unmarshal 后为float
					if dbCol.LangType == "int" {
						titem, ok := dbCol.Default.(float64)
						if ok {
							dbData[k] = int64(titem)
						} else {
							dbData[k] = 0
						}
					} else {
						dbData[k] = dbCol.Default
					}
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
		}
	} else {
		for k, dbCol := range columMap {
			//去掉主键
			if !dbCol.Key {
				_, upOk := c.Columns[k]
				if upOk {
					formVal, formOk := formData[k]
					if formOk {
						dbData[k] = formVal
					} else {
						if dbCol.LangType == "string" {
							dbData[k] = ""
						} else if dbCol.LangType == "int" || dbCol.LangType == "float" {
							dbData[k] = 0
						}
					}
				}

			}
		}
	}
	return formData, dbData, ValidResp{
		Valid: true,
		Msg:   "",
	}
}

type FormDel struct {
	Physic bool `json:"physic"`
	FormQuery
}
type FormQuery struct {
	Select []string `json:"select"`
	Where  []FormOp `json:"where"`
	Order  string   `json:"order"`
}
type FormUpdate struct {
	FormInfo
	FormQuery
}

type FormUnique struct {
	Columns []string `json:"columns"`
	Tip     string   `json:"tip"`
}

func (c *FormQuery) UnStrictParse(columMap map[string]ColumnInfo, ctx *gin.Context) bool {
	for idx, _ := range c.Where {
		item := c.Where[idx]
		val := ctx.PostForm(item.Name)
		dbCol, dbOk := columMap[item.Name]
		if dbOk {
			if val != "" {
				if dbCol.LangType == "int" {
					pval, per := strconv.ParseInt(val, 10, 64)
					if per != nil {
						pval = 0
					}
					c.Where[idx].Val = pval
				} else if dbCol.LangType == "float" {
					pval, per := strconv.ParseFloat(val, 64)
					if per != nil {
						pval = 0
					}
					c.Where[idx].Val = pval
				} else {
					c.Where[idx].Val = val
				}
			} else {
				if item.Default != nil {
					//处理json string 里面是int 但unmarshal 后为float
					if dbCol.LangType == "int" {
						titem, ok := item.Default.(float64)
						if ok {
							c.Where[idx].Val = int64(titem)
						} else {
							c.Where[idx].Val = 0
						}
					} else {
						c.Where[idx].Val = item.Default
					}
				} else {
					c.Where[idx].Val = ""
				}
			}

		} else {
			return false
		}
	}
	return true
}
func (c *FormQuery) StrictParse(columMap map[string]ColumnInfo, ctx *gin.Context) bool {
	for idx, _ := range c.Where {
		item := c.Where[idx]
		val := ctx.PostForm(item.Name)
		dbCol, dbOk := columMap[item.Name]
		if dbOk {
			if val != "" {
				if dbCol.LangType == "int" {
					pval, per := strconv.ParseInt(val, 10, 64)
					if per != nil {
						return false
					}
					c.Where[idx].Val = pval
				} else if dbCol.LangType == "float" {
					pval, per := strconv.ParseFloat(val, 64)
					if per != nil {
						return false
					}
					c.Where[idx].Val = pval
				} else {
					c.Where[idx].Val = val
				}
			} else {
				if item.Default != nil {
					//处理json string 里面是int 但unmarshal 后为float
					if dbCol.LangType == "int" {
						titem, ok := item.Default.(float64)
						if ok {
							c.Where[idx].Val = int64(titem)
						} else {
							c.Where[idx].Val = 0
						}
					} else {
						c.Where[idx].Val = item.Default
					}
				} else {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

type FormOp struct {
	Name    string      `json:"name"`
	Op      string      `json:"op"`
	On      string      `json:"on"`
	Val     interface{} `json:"-"`
	Default interface{} `json:"default"`
}

func (c *FormOp) ExpOn() bool {
	expFlag := false
	if c.On != "" {
		env := map[string]interface{}{
			c.Name: c.Val,
		}
		program, err := expr.Compile(c.On, expr.Env(env), expr.AsBool())
		if err == nil {
			output, rerr := expr.Run(program, env)
			if rerr == nil {
				expFlag = output.(bool)
			}
		}
	} else {
		expFlag = true
	}
	return expFlag
}
