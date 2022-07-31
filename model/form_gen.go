package models

import (
	"encoding/json"
	"firefly/utils"
	"io/ioutil"
	"path/filepath"
)

func ModLoadFile(filePath string) ModInfo {
	if utils.PathExists(filePath) {
		_, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			return x.(ModInfo)
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return ModInfo{}
		}
		info := ModLoadByte(data)
		utils.SetCache(fileName, info)
		return info
	} else {
		return ModInfo{}
	}
}
func ModLoadByte(jsonByte []byte) ModInfo {
	var info ModInfo
	json.Unmarshal(jsonByte, &info)
	return info
}

func FormAddLoadFile(filePath string) FormAdd {
	if utils.PathExists(filePath) {
		_, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			return x.(FormAdd)
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return FormAdd{}
		}
		info := FormAddLoadByte(data)
		utils.SetCache(fileName, info)
		return info
	} else {
		return FormAdd{}
	}
}
func FormAddLoadByte(jsonByte []byte) FormAdd {
	var info FormAdd
	json.Unmarshal(jsonByte, &info)
	return info
}

func FormUpdateLoadFile(filePath string) FormUpdate {
	if utils.PathExists(filePath) {
		_, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			return x.(FormUpdate)
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return FormUpdate{}
		}
		info := FormUpdateLoadByte(data)
		utils.SetCache(fileName, info)
		return info
	} else {
		return FormUpdate{}
	}
}
func FormUpdateLoadByte(jsonByte []byte) FormUpdate {
	var info FormUpdate
	json.Unmarshal(jsonByte, &info)
	return info
}

func FormQueryLoadFile(filePath string) FormQuery {
	if utils.PathExists(filePath) {
		_, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			return x.(FormQuery)
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return FormQuery{}
		}
		info := FormQueryLoadByte(data)
		utils.SetCache(fileName, info)
		return info
	} else {
		return FormQuery{}
	}
}
func FormQueryLoadByte(jsonByte []byte) FormQuery {
	var info FormQuery
	uer := json.Unmarshal(jsonByte, &info)
	if uer != nil {
		println(uer.Error())
	}
	return info
}
