package spider

import (
	"movie/configs"
	"movie/services/spider/zyw"
)

// Spider 定义一个采集接口
type Spider interface {
	// SpiderSearch 接口需要实现的方法列表
	SpiderSearch(site configs.SpiderSites, kw string) (data []Movie, msg string)
}

// ExternalApi 外部接口回调数据处理函数映射表
var ExternalApi = map[string]ExternalSpider{
	"zyw":   &zyw.ZYW{},
	"xlzyw": &zyw.XlZyw{},
	"gszyw": &zyw.GSZyw{},
}

// ExternalSpider todo 外部数据调用接口，接口统一需实现的方法如下
type ExternalSpider interface {
	ExternalGetReport(configs.SpiderSites, string) (zyw.ZySearchResult, error) // 主动查询三方结果
	ExternalGetById(configs.SpiderSites, string) (zyw.ZyDetailResult, error)   // 根据影片ID查询三方数据
}
