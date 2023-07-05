package easycache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var C *cache.Cache

func init() {
	// 实例化缓存组件
	C = cache.New(3*60*time.Minute, 3*60*time.Minute)
}
