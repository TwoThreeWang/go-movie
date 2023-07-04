package easycache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var C *cache.Cache

func init() {
	// 实例化缓存组件
	C = cache.New(5*time.Minute, 10*time.Minute)
}
