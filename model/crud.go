package models

import (
	"encoding/json"
	"firefly/utils"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type CrudEngine struct {
	Name    string
	Process string
	Node    string
}

func NewCrudEngine(code string) CrudEngine {
	codeList := strings.Split(code, ".")
	return CrudEngine{
		Name:    codeList[0] + ".json",
		Process: strings.ToLower(codeList[1]),
		Node:    codeList[2],
	}
}

type CrudInfo struct {
	Mod    ModInfo               `json:"mod"`
	List   map[string]FormList   `json:"list"`
	Page   map[string]FormPage   `json:"page"`
	Get    map[string]FormGet    `json:"get"`
	Add    map[string]FormAdd    `json:"add"`
	Update map[string]FormUpdate `json:"update"`
	Del    map[string]FormDel    `json:"del"`
}

func (c *CrudInfo) checkDisable() {
	for _, v := range c.Update {
		v.checkDisable()
	}
	for _, v := range c.List {
		v.checkDisable()
	}
	for _, v := range c.Del {
		v.checkDisable()
	}
	for _, v := range c.Add {
		v.checkDisable()
	}
	for _, v := range c.Page {
		v.checkDisable()
	}
	for _, v := range c.Get {
		v.checkDisable()
	}
}

func LoadCrudFile(name string) (CrudInfo, bool) {
	filePath := "resource/crud/" + name
	if !utils.PathExists(filePath) {
		return CrudInfo{}, false
	}
	_, fileName := filepath.Split(filePath)
	if x, found := utils.GetCache(fileName); found {
		return x.(CrudInfo), true
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return CrudInfo{}, false
	}
	info := LoadCrudByte(data)
	utils.SetCache(fileName, info)
	return info, true
}
func LoadCrudByte(jsonByte []byte) CrudInfo {
	var info CrudInfo
	json.Unmarshal(jsonByte, &info)
	info.checkDisable()
	return info
}
