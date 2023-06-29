package spider

import (
	"encoding/json"
	"fmt"
	"movie/configs"
	"movie/controllers/admin"
	"movie/utils/easyhttp"
	"movie/utils/easylog"
	"strconv"
	"time"
)

func doSearch(site configs.SpiderSites, kw string) (movie []Movie, err error) {
	// 获取当前时间的13位时间戳
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	url := site.BaseUrl + "/index.php/ajax/suggest?mid=1&wd=" + kw + "&limit=10&imestamp=" + strconv.FormatInt(timestamp, 10)
	// Get请求接口
	jsonStr, err := easyhttp.Get(url, nil, nil)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	// 解析JSON响应
	var jsondata MacCmsSearchResult
	err = json.Unmarshal(jsonStr, &jsondata)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	movie = jsondata.List
	return
}

// MacCmsSearch 苹果CMS 搜索接口
func MacCmsSearch(baseUrl string, kw string) (data []Movie, msg string) {
	var sites []configs.SpiderSites
	sites, err := admin.GetSites()
	if err != nil {
		msg = err.Error()
		return
	}
	for _, site := range sites {
		movies, err := doSearch(site, kw)
		if err != nil {
			msg = err.Error()
			return
		}
		// 遍历JSON数组，增加数据来源和播放链接
		for _, movie := range movies {
			movie.Source = baseUrl
			movie.VideoUrl = fmt.Sprintf(site.PlayUrl, baseUrl, strconv.Itoa(movie.Id))
			data = append(data, movie)
		}
	}
	return
}
