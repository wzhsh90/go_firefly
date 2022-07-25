package models

import (
	"encoding/json"
	"firefly/utils"
	"io/ioutil"
	"path/filepath"
)

type LoginUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	Id      string `db:"id" form:"id" json:"id"`
	ComDesc string `db:"com_desc" form:"com_desc" json:"com_desc"`
	ComName string `db:"com_name" form:"com_name" json:"com_name"`
}

type ModInfo struct {
	Name    string               `json:"name"`
	Table   ModTable             `json:"table"`
	Columns map[string]ModColumn `json:"columns"`
}

func (c *ModInfo) GetKeyCols() map[string]ModColumn {
	keys := make(map[string]ModColumn)
	for k, v := range c.Columns {
		if v.Key {
			keys[k] = v
		}
	}
	return keys
}

type ModTable struct {
	Name   string `json:"name"`
	ZhName string `json:"zh_name"`
}
type ModColumn struct {
	ZhName   string `json:"zh_name"`
	DbType   string `json:"db_type"`
	LangType string `json:"lang_type"`
	Key      bool   `json:"key"`
	Index    string `json:"index"`
}

type FormAddInfo struct {
	Exits struct {
		Columns []string `json:"columns"`
		Tip     string   `json:"tip"`
	} `json:"exits"`
	Columns map[string]FormAddColumn `json:"columns"`
}

func (c *FormAddInfo) GetRules() map[string]interface{} {
	rules := make(map[string]interface{})
	for k, v := range c.Columns {
		if v.Rule != "" {
			rules[k] = v.Rule
		}
	}
	return rules
}

type FormAddColumn struct {
	ZhName string `json:"zh_name"`
	Dom    string `json:"dom"`
	Rule   string `json:"rule"`
}

func ModInfoGenFile(filePath string) ModInfo {
	if utils.PathExists(filePath) {
		_, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			return x.(ModInfo)
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return ModInfo{}
		}
		info := ModInfoGen(data)
		utils.SetCache(fileName, info)
		return info
	} else {
		return ModInfo{}
	}
}
func ModInfoGen(jsonByte []byte) ModInfo {
	var info ModInfo
	json.Unmarshal(jsonByte, &info)
	return info
}

func FormAddModGenFile(filePath string) FormAddInfo {
	if utils.PathExists(filePath) {
		_, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			return x.(FormAddInfo)
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return FormAddInfo{}
		}
		info := FormAddInfoGen(data)
		utils.SetCache(fileName, info)
		return info
	} else {
		return FormAddInfo{}
	}
}
func FormAddInfoGen(jsonByte []byte) FormAddInfo {
	var info FormAddInfo
	json.Unmarshal(jsonByte, &info)
	return info
}
