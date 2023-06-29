package spider

import "movie/configs"

// Spider 采集接口
type Spider interface {
	ApiSearch(site configs.SpiderSites, kw string) (data []Movie, msg string)
}
