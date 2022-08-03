package models

import (
	"encoding/json"
	"firefly/utils"
	"io/ioutil"
	"path/filepath"
)

type CurdInfo struct {
	Mod    ModInfo    `json:"mod"`
	List   FormList   `json:"list"`
	Add    FormAdd    `json:"add"`
	Update FormUpdate `json:"update"`
	Del    FormDel    `json:"del"`
}

func (c *CurdInfo) IsEnable() {
	c.Del.IsEnable()
	c.List.IsEnable()
	c.Add.IsEnable()
	c.Update.IsEnable()
}

func LoadCrudFile(filePath string) CurdInfo {
	if utils.PathExists(filePath) {
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
	} else {
		return CurdInfo{}
	}
}
func LoadCrudByte(jsonByte []byte) CurdInfo {
	var info CurdInfo
	json.Unmarshal(jsonByte, &info)
	info.IsEnable()
	return info
}
