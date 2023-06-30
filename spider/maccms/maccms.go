package maccms

import (
	"encoding/json"
	"fmt"
	"movie/configs"
	"movie/spider"
	"movie/utils/easyhttp"
	"movie/utils/easylog"
	"strconv"
	"time"
)

// MacCmsSearch 苹果CMS搜索接口实现
type MacCmsSearch struct{}

func (m MacCmsSearch) SpiderSearch(site configs.SpiderSites, kw string) (data []spider.Movie, msg string) {
	movies, err := doSearch(site, kw)
	if err != nil {
		msg = err.Error()
		return
	}
	// 遍历JSON数组，增加数据来源和播放链接
	for _, movie := range movies {
		movie.Source = site.BaseUrl
		movie.VideoUrl = fmt.Sprintf(site.PlayUrl, site.BaseUrl, strconv.Itoa(movie.Id))
		data = append(data, movie)
	}
	return
}

// 实际采集逻辑
func doSearch(site configs.SpiderSites, kw string) (movie []spider.Movie, err error) {
	// 获取当前时间的13位时间戳
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	url := site.BaseUrl + "/index.php/ajax/suggest"
	param := map[string]interface{}{
		"mid":      "1",
		"wd":       kw,
		"limit":    "10",
		"imestamp": strconv.FormatInt(timestamp, 10),
	}
	header := map[string]string{
		"Content-Type": "application/json",
	}
	// Get请求接口
	jsonStr, err := easyhttp.Post(url, header, param)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	fmt.Println(url + " 有结果")
	// 解析JSON响应
	var jsondata spider.MacCmsSearchResult
	err = json.Unmarshal(jsonStr, &jsondata)
	if err != nil {
		fmt.Println(url + " 错误结果 " + string(jsonStr))
		easylog.Log.Error(err)
		return
	}
	movie = jsondata.List
	return
}
