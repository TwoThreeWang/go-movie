package zyw

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"movie/configs"
	"movie/utils/easyhttp"
	"movie/utils/easylog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type GSZyw struct{}

// ExternalGetReport 主动查询三方结果
func (zyw *GSZyw) ExternalGetReport(site configs.SpiderSites, kw string) (movie ZySearchResult, err error) {
	url := site.WebUrl + "?wd=" + url.QueryEscape(kw)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		easylog.Log.Error(err)
		return
	}

	items := htmlquery.Find(doc, "//td[@class='yp']")
	for _, item := range items {
		aTag := htmlquery.FindOne(item, ".//a")
		link := htmlquery.SelectAttr(aTag, "href")
		// 使用正则表达式匹配ID
		re := regexp.MustCompile(`id/(\d+)\.html`)
		matches := re.FindStringSubmatch(link)
		if len(matches) > 1 {
			id := matches[1]
			movieDetail, err := zyw.ExternalGetById(site, id)
			if err != nil {
				easylog.Log.Error(err)
			}
			if len(movieDetail.List) > 0 {
				movie.List = append(movie.List, movieDetail.List[0])
			}
		}
	}
	movie.Source = site.Key
	return
}

// ExternalGetById 根据影片ID查询三方数据
func (zyw *GSZyw) ExternalGetById(site configs.SpiderSites, vid string) (movie ZyDetailResult, err error) {
	url := site.BaseUrl + "/provide/vod/?ac=detail&ids=" + vid
	param := map[string]interface{}{}
	header := map[string]string{
		"Content-Type": "application/json",
	}
	// Get请求接口
	jsonStr, err := easyhttp.Get(url, header, param)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	// 解析JSON响应
	err = json.Unmarshal(jsonStr, &movie)
	if err != nil {
		fmt.Println(url + " 错误结果 " + string(jsonStr))
		easylog.Log.Error(err)
		return
	}
	movie.Source = site.Key
	return
}
