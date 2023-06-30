package spider

import (
	"movie/configs"
)

// Spider 定义一个采集接口
type Spider interface {
	// 接口需要实现的方法列表
	SpiderSearch(site configs.SpiderSites, kw string) (data []Movie, msg string)
}
