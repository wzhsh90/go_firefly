package models

import (
	"firefly/utils"
	"github.com/antonmedv/expr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
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

func parseDefaultVal(langType string, defVal interface{}) interface{} {

	//处理json string 里面是int 但unmarshal 后为float
	if defVal != nil {
		if langType != "int" {
			return defVal
		}
		item, ok := defVal.(float64)
		if ok {
			return int64(item)
		} else {
			return 0
		}
	} else {
		if langType == "int" || langType == "float" {
			return 0
		} else {
			return ""
		}
	}

}
func (c *ColumnInfo) getInitVal() interface{} {
	if c.LangType == "int" || c.LangType == "float" {
		return 0
	} else {
		return ""
	}
}
func (c *ColumnInfo) parseDefaultVal() interface{} {
	return parseDefaultVal(c.LangType, c.Default)
}
func (c *ColumnInfo) checkVal(val string) CheckColumnResp {

	if c.LangType == "int" {
		pval, per := strconv.ParseInt(val, 10, 64)
		msg := ""
		if per != nil {
			pval = 0
			msg = c.ZhName + "数据不合法"
		}
		return CheckColumnResp{
			Valid:    per == nil,
			Msg:      msg,
			Default:  0,
			ParseVal: pval,
		}
	} else if c.LangType == "float" {
		pval, per := strconv.ParseFloat(val, 64)
		msg := ""
		if per != nil {
			pval = 0
			msg = c.ZhName + "数据不合法"
		}
		return CheckColumnResp{
			Valid:    per == nil,
			Msg:      msg,
			Default:  0,
			ParseVal: pval,
		}
	} else {
		return CheckColumnResp{
			Valid:    true,
			Msg:      "",
			Default:  "",
			ParseVal: val,
		}
	}
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
type CheckColumnResp struct {
	Valid    bool
	Msg      string
	Default  interface{}
	ParseVal interface{}
}

func (c *FormInfo) GetFormData(columMap map[string]ColumnInfo, ctx *gin.Context, create bool) (map[string]interface{}, map[string]interface{}, ValidResp) {
	formData := make(map[string]interface{})
	dbData := make(map[string]interface{})
	//生成form data
	for k, _ := range c.Columns {
		val := ctx.PostForm(k)
		dbCol, dbOk := columMap[k]
		if !dbOk {
			return nil, nil, ValidResp{
				Valid: false,
				Msg:   "数据不合法,存在未定义字段" + k,
			}
		}
		if len(val) > 0 {
			checkResp := dbCol.checkVal(val)
			if !checkResp.Valid {
				return nil, nil, ValidResp{
					Valid: false,
					Msg:   checkResp.Msg,
				}
			} else {
				formData[k] = checkResp.ParseVal
			}
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
				dbData[k] = dbCol.parseDefaultVal()
				//判断是key并且值为空，则默认创建id
				if dbCol.Key && dbCol.LangType == "string" && dbData[k] == "" {
					dbData[k] = utils.NewId()
				}
			}
		}
	} else {
		for k, dbCol := range columMap {
			//去掉主键
			if dbCol.Key {
				continue
			}
			_, upOk := c.Columns[k]
			if !upOk {
				continue
			}
			formVal, formOk := formData[k]
			if formOk {
				dbData[k] = formVal
			} else {
				dbData[k] = dbCol.getInitVal()
			}
		}
	}
	return formData, dbData, ValidResp{
		Valid: true,
		Msg:   "",
	}
}

type FormAdd struct {
	Opt
	FormInfo
}

func (c *FormAdd) checkDisable() {
	if len(c.Columns) <= 0 {
		c.Opt.Disable = true
	}
}
func (c *FormGet) checkDisable() {
	if len(c.Select) <= 0 {
		c.Opt.Disable = true
	} else {
		c.whereParse()
	}
}
func (c *FormPage) checkDisable() {
	if len(c.Select) <= 0 {
		c.Opt.Disable = true
	} else {
		c.whereParse()
	}
}

type FormList struct {
	From  string `json:"from"`
	Order string `json:"order"`
	FormQuery
}
type FormPage struct {
	From  string `json:"from"`
	Order string `json:"order"`
	FormQuery
}

type FormGet struct {
	From string `json:"from"`
	FormQuery
}

func (c *FormList) checkDisable() {
	if len(c.Select) <= 0 {
		c.Opt.Disable = true
	} else {
		c.whereParse()
	}
}

type FormDel struct {
	Fake bool `json:"fake"`
	FormQuery
}
type FormQuery struct {
	Opt
	Select []string `json:"select"`
	Where  []FormOp `json:"where"`
}

func (c *FormQuery) whereParse() {
	for idx, _ := range c.Where {
		c.Where[idx].FormatName()
	}
}
func (c *FormQuery) UnStrictParse(columMap map[string]ColumnInfo, ctx *gin.Context) bool {
	for idx, _ := range c.Where {
		item := c.Where[idx]
		val := ctx.PostForm(item.PlainName)
		dbCol, dbOk := columMap[item.PlainName]
		if !dbOk {
			return false
		}
		if val != "" {
			checkResp := dbCol.checkVal(val)
			if checkResp.Valid {
				c.Where[idx].Val = checkResp.ParseVal
			} else {
				c.Where[idx].Val = checkResp.Default
			}
		} else {
			c.Where[idx].Val = item.parseDefaultVal(dbCol.LangType)
		}
	}
	return true
}
func (c *FormQuery) StrictParse(columMap map[string]ColumnInfo, ctx *gin.Context) bool {
	for idx, _ := range c.Where {
		item := c.Where[idx]
		val := ctx.PostForm(item.PlainName)
		dbCol, dbOk := columMap[item.PlainName]
		if !dbOk {
			return false
		}
		if val != "" {
			checkResp := dbCol.checkVal(val)
			if !checkResp.Valid {
				return false
			} else {
				c.Where[idx].Val = checkResp.ParseVal
			}
		} else {
			if item.Default != nil {
				c.Where[idx].Val = item.parseDefaultVal(dbCol.LangType)
			} else {
				return false
			}
		}
	}
	return true
}

type FormUpdate struct {
	FormInfo
	FormQuery
}

func (c *FormUpdate) checkDisable() {
	if len(c.Where) <= 0 {
		c.Opt.Disable = true
	} else {
		c.whereParse()
	}
}

type Opt struct {
	Disable bool `json:"disable"`
}

func (c *FormDel) checkDisable() {
	if len(c.Where) <= 0 {
		c.Opt.Disable = true
	} else {
		c.whereParse()
	}
}

type FormUnique struct {
	Columns []string `json:"columns"`
	Tip     string   `json:"tip"`
}

type FormOp struct {
	Name      string      `json:"name"`
	Prefix    string      `json:"-"`
	PlainName string      `json:"-"`
	Op        string      `json:"op"`
	On        string      `json:"on"`
	Val       interface{} `json:"-"`
	Default   interface{} `json:"default"`
}

func (c *FormOp) parseDefaultVal(langType string) interface{} {
	return parseDefaultVal(langType, c.Default)
}
func (c *FormOp) FormatName() {
	if strings.Contains(c.Name, ".") {
		nameList := strings.Split(c.Name, ".")
		c.Prefix = nameList[0]
		lastName := nameList[len(nameList)-1]
		c.PlainName = lastName
	} else {
		c.PlainName = c.Name
	}
}
func (c *FormOp) ExpOn() bool {
	if c.On == "" {
		return true
	}
	env := make(map[string]interface{})
	env[c.PlainName] = c.Val
	if c.Prefix != "" {
		env[c.Prefix] = map[string]interface{}{
			c.PlainName: c.Val,
		}
	}
	program, err := expr.Compile(c.On, expr.Env(env), expr.AsBool())
	if err == nil {
		output, rerr := expr.Run(program, env)
		if rerr == nil {
			return output.(bool)
		}
	}
	return false
}
