package spider

// MacCmsSearchResult 苹果CMS搜索接口返回结构体
type MacCmsSearchResult struct {
	Code      int     `json:"code"`
	Msg       string  `json:"msg"`
	Page      int     `json:"page"`
	Pagecount int     `json:"pagecount"`
	Limit     int     `json:"limit"`
	Total     int     `json:"total"`
	List      []Movie `json:"list"` // 电影结果列表
	Url       string  `json:"url"`
}

type Movie struct {
	Id       int    `json:"id"`       // 电影ID
	Name     string `json:"name"`     // 电影名称
	En       string `json:"en"`       // 电影英文名称
	Pic      string `json:"pic"`      // 封面图
	Source   string `json:"source"`   // 来源
	VideoUrl string `json:"videoUrl"` // 观看链接
}
