package douban

import (
	"encoding/json"
	"errors"
	"fmt"
	"movie/utils/easyhttp"
	"movie/utils/easylog"
)

func DoubanHot(movie_type string) (data ResultDoubanHot, err error) {
	var url string
	if movie_type == "movie" {
		url = "https://movie.douban.com/j/search_subjects?type=movie&tag=热门&page_limit=50&page_start=0"
	} else if movie_type == "tv" {
		url = "https://movie.douban.com/j/search_subjects?type=tv&tag=热门&page_limit=50&page_start=0"
	} else if movie_type == "show" {
		url = "https://movie.douban.com/j/search_subjects?type=tv&tag=综艺&page_limit=50&page_start=0"
	} else if movie_type == "cartoon" {
		url = "https://movie.douban.com/j/search_subjects?type=tv&tag=日本动画&page_limit=50&page_start=0"
	} else {
		err = errors.New("参数错误")
		return
	}
	header := map[string]string{
		"Referer":    "https://movie.douban.com/",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67",
		"Host":       "movie.douban.com",
		"Cookie":     "bid=jyRD6omH2o0; douban-fav-remind=1; ll=\"108288\"; ct=y; frodotk_db=\"c09da24285e2c1b041d8e6c791a84510\"; push_noty_num=0; push_doumail_num=0; ap_v=0,6.0; bid=RsvEZPEuiWQ",
		"Accept":     "*/*",
		//"Accept-Encoding":    "gzip, deflate, br",
		"Accept-Language":    "zh-CN,zh;q=0.9,en;q=0.8",
		"Connection":         "keep-alive",
		"Dnt":                "1",
		"Sec-Ch-Ua":          "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", Microsoft Edge\";v=\"114\"",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": "\"Windows\"",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"X-Requested-With":   "XMLHttpRequest",
	}
	// Get请求接口
	jsonStr, err := easyhttp.Get(url, header, nil)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	// 解析JSON响应
	err = json.Unmarshal(jsonStr, &data)
	if err != nil {
		fmt.Println(url + " 错误结果 " + string(jsonStr))
		easylog.Log.Error(err)
		return
	}
	return
}

type ResultDoubanHot struct {
	Subjects []struct {
		EpisodesInfo string `json:"episodes_info"`
		Rate         string `json:"rate"`
		CoverX       int    `json:"cover_x"`
		Title        string `json:"title"`
		Url          string `json:"url"`
		Playable     bool   `json:"playable"`
		Cover        string `json:"cover"`
		Id           string `json:"id"`
		CoverY       int    `json:"cover_y"`
		IsNew        bool   `json:"is_new"`
	} `json:"subjects"`
}
