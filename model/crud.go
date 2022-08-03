package models

import (
	"encoding/json"
	"firefly/utils"
	"io/ioutil"
	"path/filepath"
)

type CurdInfo struct {
	Mod    ModInfo               `json:"mod"`
	List   map[string]FormList   `json:"list"`
	Add    map[string]FormAdd    `json:"add"`
	Update map[string]FormUpdate `json:"update"`
	Del    map[string]FormDel    `json:"del"`
}

func (c *CurdInfo) checkDisable() {
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
}

func LoadCrudFile(filePath string) CurdInfo {
	if !utils.PathExists(filePath) {
		return CurdInfo{}
	}
	_, fileName := filepath.Split(filePath)
	if x, found := utils.GetCache(fileName); found {
		return x.(CurdInfo)
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return CurdInfo{}
	}
	info := LoadCrudByte(data)
	utils.SetCache(fileName, info)
	return info
}
func LoadCrudByte(jsonByte []byte) CurdInfo {
	var info CurdInfo
	json.Unmarshal(jsonByte, &info)
	info.checkDisable()
	return info
}
