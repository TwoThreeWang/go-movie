package configs

type SpiderSites struct {
	Key       string `json:"key"`       // 网站简称
	BaseUrl   string `json:"baseurl"`   // 基础url
	PlayUrl   string `json:"playurl"`   // 播放链接
	SpiderKey string `json:"spiderkey"` // 采集函数标识
}
