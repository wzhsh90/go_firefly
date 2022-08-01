package models

import (
	"encoding/json"
	"firefly/utils"
	"io/ioutil"
	"path"
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

func FormAddLoadFile(filePath string) (FormAdd, ModInfo) {
	if utils.PathExists(filePath) {
		dir, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			cacheInfo := x.(FormAdd)
			return cacheInfo, ModLoadFile(path.Join(cacheInfo.Dir, cacheInfo.Mod))
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return FormAdd{}, ModInfo{}
		}
		info := FormAddLoadByte(data)
		info.Dir = dir
		utils.SetCache(fileName, info)
		return info, ModLoadFile(path.Join(info.Dir, info.Mod))
	} else {
		return FormAdd{}, ModInfo{}
	}
}
func FormAddLoadByte(jsonByte []byte) FormAdd {
	var info FormAdd
	json.Unmarshal(jsonByte, &info)
	return info
}

func FormUpdateLoadFile(filePath string) (FormUpdate, ModInfo) {
	if utils.PathExists(filePath) {
		dir, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			cacheInfo := x.(FormUpdate)
			return cacheInfo, ModLoadFile(path.Join(cacheInfo.Dir, cacheInfo.Mod))
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return FormUpdate{}, ModInfo{}
		}
		info := FormUpdateLoadByte(data)
		info.Dir = dir
		utils.SetCache(fileName, info)
		return info, ModLoadFile(path.Join(info.Dir, info.Mod))
	} else {
		return FormUpdate{}, ModInfo{}
	}
}
func FormUpdateLoadByte(jsonByte []byte) FormUpdate {
	var info FormUpdate
	json.Unmarshal(jsonByte, &info)
	return info
}

func FormQueryLoadFile(filePath string) (FormQuery, ModInfo) {
	if utils.PathExists(filePath) {
		dir, fileName := filepath.Split(filePath)
		if x, found := utils.GetCache(fileName); found {
			cacheInfo := x.(FormQuery)
			return cacheInfo, ModLoadFile(path.Join(cacheInfo.Dir, cacheInfo.Mod))
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return FormQuery{}, ModInfo{}
		}
		info := FormQueryLoadByte(data)
		info.Dir = dir
		utils.SetCache(fileName, info)
		return info, ModLoadFile(path.Join(info.Dir, info.Mod))
	} else {
		return FormQuery{}, ModInfo{}
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
