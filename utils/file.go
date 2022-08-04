package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// DirNotExists 校验目录是否存在
func DirNotExists(dir string) bool {
	dir = strings.TrimPrefix(dir, "fs://")
	dir = strings.TrimPrefix(dir, "file://")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return true
	}
	return false
}

// DirAbs 文件绝对路径
func DirAbs(dir string) string {
	dir = strings.TrimPrefix(dir, "fs://")
	dir = strings.TrimPrefix(dir, "file://")
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		log.Panic("获取绝对路径错误 %s %s", dir, err)
	}
	return dirAbs
}
