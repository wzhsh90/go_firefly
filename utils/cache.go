package utils

import (
	"github.com/patrickmn/go-cache"
)

var cacheBox *cache.Cache

func init() {
	//cacheBox = cache.New(5*time.Minute, 10*time.Minute)
	cacheBox = cache.New(0, 0)
}
func SetCache(key string, value interface{}) {
	cacheBox.Set(key, value, cache.NoExpiration)
}
func GetCache(key string) (interface{}, bool) {
	return cacheBox.Get(key)
}
func ExistCache(key string) bool {
	_, flag := cacheBox.Get(key)
	return flag
}
