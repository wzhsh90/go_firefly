package middleware

import (
	models "firefly/model"
	"firefly/utils"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

func init() {
	WatchModel(models.ResourcePath)
}

// WatchModel 监听业务接口更新
func WatchModel(root string) {
	if utils.DirNotExists(root) {
		return
	}
	root = utils.DirAbs(root)
	go Watch(root, func(op string, filename string) {
		if !strings.HasSuffix(filename, ".json") {
			return
		}
		if op == "write" || op == "create" {
			key := SpecName(root, filename)
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Println("nt exist")
				return
			}
			utils.SetCache(key, models.LoadCrudByte(data))
			fmt.Println(color.GreenString("file %s %s", key, op))
		} else if op == "remove" || op == "rename" {
			key := SpecName(root, filename)
			if utils.ExistCache(key) {
				utils.DeleteCache(key)
				fmt.Println(color.GreenString("file %s %s", key, op))
			}
		}
	})
}
func SpecName(root string, file string) string {
	filename := strings.TrimPrefix(file, root+"/")
	namer := strings.Split(filename, ".")
	return namer[0]
}

var watchGroup sync.WaitGroup
var watchOp = map[fsnotify.Op]string{
	fsnotify.Create: "create",
	fsnotify.Write:  "write",
	fsnotify.Remove: "remove",
	fsnotify.Rename: "rename",
	fsnotify.Chmod:  "chmod",
}

// Watch 监听目录
func Watch(root string, cb func(op string, file string)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		watcher.Close()
		watchGroup.Done()
	}()

	watchGroup.Add(1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 监听子目录
				if event.Op == fsnotify.Create {
					file, err := os.Open(event.Name)
					if err == nil {
						fi, err := file.Stat()
						file.Close()
						if err == nil && fi.IsDir() {
							Watch(event.Name, cb)
						}
					}
				}

				cb(watchOp[event.Op], event.Name)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(color.GreenString("Watching: %s", root))

	// 监听子目录
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("not exit")
			return err
		}

		if path == root {
			return nil
		}

		if d.IsDir() {
			go Watch(path, cb)
		}
		return nil
	})

	watchGroup.Wait()

}
