package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"log"
	"movie/configs"
	"movie/repositories"
	"movie/services/spider"
	"movie/services/spider/douban"
	"movie/services/spider/zyw"
	"movie/utils/easycache"
	"movie/utils/easylog"
	"movie/utils/result"
	"movie/utils/try"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetSites 获取已有站点数据
func GetSites() (sites []configs.SpiderSites, err error) {
	err = viper.UnmarshalKey("sites", &sites)
	return
}

// Search 搜索接口
func Search(c *gin.Context) {
	kw := c.Query("kw") // 搜索关键词
	if kw != "" {
		if x, found := easycache.C.Get(kw); found {
			datas := x.([]zyw.ZyMovie)
			result.Ok(c, datas)
			c.Abort()
			return
		}
		try.Try(func() {
			sites, err := GetSites()
			if err != nil {
				result.Fail(c, result.ResultError, err.Error())
				c.Abort()
				return
			}
			// 创建结果通道
			resultChan := make(chan zyw.ZySearchResult, len(sites))
			// 创建超时通道
			timeOut, err := strconv.Atoi(viper.GetString("spider_timeout"))
			if err != nil {
				result.Fail(c, result.ResultError, err.Error())
				c.Abort()
				return
			}
			timeoutChan := time.After(time.Duration(timeOut) * time.Second)
			// 经过context的WithTimeout设置一个有效时间为800毫秒的context
			// 该context会在耗尽800毫秒后或者方法执行完成后结束,结束的时候会向通道ctx.Done发送信号
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
			// 注意,这里要记得调用cancel(),否则即便提早执行完了,还要傻傻等到800毫秒后context才会被释放
			defer cancel()
			var wg sync.WaitGroup
			for _, site := range sites {
				site := site
				wg.Add(1)
				go func(ctx context.Context) {
					defer wg.Done()
					f, ok := spider.ExternalApi[site.SpiderKey]
					if !ok {
						easylog.Log.Error(site.SpiderKey + " 没找到映射采集接口")
					} else {
						// 动态调用采集函数
						dataNew, err := f.ExternalGetReport(site, kw)
						if err != nil {
							easylog.Log.Error(err.Error())
							return
						}
						// 将响应结果发送到结果通道
						resultChan <- dataNew
					}
				}(ctx)
			}
			// 等待所有响应结果返回或者超时
			var datas []zyw.ZyMovie
			for i := 0; i < len(sites); i++ {
				select {
				case res := <-resultChan:
					// 返回结果处理
					source := res.Source
					for _, data := range res.List {
						data.Source = source
						data.VodContent = strip.StripTags(data.VodContent)
						datas = append(datas, data)
					}
				case <-timeoutChan:
					cancel()
					fmt.Println("timeout!!!")
					if len(datas) > 0 {
						easycache.C.Set(kw, datas, cache.DefaultExpiration)
						easylog.Log.Info("开始存库和写入搜索历史")
						go saveDb(datas) // 结果保存到数据库
						go addSearchHistory(kw)
					}
					result.Ok(c, datas)
					c.Abort()
					return
				}
			}
			wg.Wait()
			if len(datas) > 0 {
				easycache.C.Set(kw, datas, cache.DefaultExpiration)
				easylog.Log.Info("开始存库和写入搜索历史")
				go saveDb(datas) // 结果保存到数据库
				go addSearchHistory(kw)
			}
			result.Ok(c, datas)
			c.Abort()
			return
		}).CatchAll(func(err error) {
			easylog.Log.Error(err.Error())
			result.Fail(c, result.ResultError, err.Error())
			c.Abort()
			return
		})
	} else {
		result.FailNoMsg(c, result.InvalidArgs)
		c.Abort()
		return
	}
}

// 影片信息保存到数据库
func saveDb(datas []zyw.ZyMovie) {
	easylog.Log.Info("存入数据库：", datas)
	db := repositories.GetDB()
	var movies repositories.Movies
	for _, movie := range datas {
		movie.UpdatedAt = time.Now()
		movies = repositories.Movies(movie)
		var mov repositories.Movies
		res := db.Where("source = ? AND vod_id = ?", movies.Source, movies.VodId).First(&mov)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			easylog.Log.Info("新建：", movie)
			// 如果没有找到记录，则创建一条新的记录
			res = db.Create(&movies)
			if res.Error != nil {
				easylog.Log.Error(res.Error)
			}
		} else if res.Error != nil {
			easylog.Log.Error(res.Error)
		} else {
			easylog.Log.Info("更新：", movie)
			// 如果找到了记录，则更新该记录
			res = db.Save(&movies)
			if res.Error != nil {
				easylog.Log.Error(res.Error)
			}
		}
	}
}

// DoubanHot 豆瓣热门电影数据获取
func DoubanHot(c *gin.Context) {
	movie_type := c.Query("type") // 热搜类型，movie：电影；tv：电视剧
	if movie_type == "" {
		result.FailNoMsg(c, result.InvalidArgs)
		c.Abort()
		return
	}
	if x, found := easycache.C.Get(movie_type); found {
		data := x.(douban.ResultDoubanHot)
		result.Ok(c, data.Subjects)
		c.Abort()
		return
	}
	data, err := douban.DoubanHot(movie_type)
	if err != nil {
		result.FailNoMsg(c, result.InvalidArgs)
		c.Abort()
		return
	}
	if len(data.Subjects) > 0 {
		easycache.C.Set(movie_type, data, 6*time.Hour)
	}
	result.Ok(c, data.Subjects)
}

// Index 接口首页
func Index(c *gin.Context) {
	var data = map[string]string{
		"app_name": "api v1.0",
	}
	result.Ok(c, data)
}

// GetMovieInfo 查询影片信息接口
func GetMovieInfo(c *gin.Context) {
	source := c.Query("source")   // 来源
	VodId := c.Query("id")        // 影片ID
	refresh := c.Query("refresh") // 强制更新
	cacheKey := source + VodId
	if refresh != "1" {
		if x, found := easycache.C.Get(cacheKey); found {
			datas := x.(repositories.Movies)
			result.Ok(c, datas)
			c.Abort()
			return
		}
	}
	db := repositories.GetDB()
	var movie repositories.Movies
	res := db.Where("source = ? and vod_id = ?", source, VodId).First(&movie)
	if res.Error != nil {
		go UpdataMovieInfo(source, VodId)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			result.FailNoMsg(c, result.NotFound)
			c.Abort()
			return
		}
		result.FailNoMsg(c, result.DatabaseError)
		c.Abort()
		return
	}
	now := time.Now()
	if movie.UpdatedAt.Year() != now.Year() || movie.UpdatedAt.Month() != now.Month() || movie.UpdatedAt.Day() != now.Day() || refresh == "1" {
		// 数据最后更新时间不是今天，后台更新这条数据
		go UpdataMovieInfo(source, VodId)
	}
	easycache.C.Set(cacheKey, movie, cache.DefaultExpiration)
	result.Ok(c, movie)
}

// UpdataMovieInfo 根据ID和来源更新影片信息
func UpdataMovieInfo(source string, VodId string) {
	sites, _ := GetSites()
	var site configs.SpiderSites
	for _, site = range sites {
		if site.Key == source {
			break
		}
	}
	f, ok := spider.ExternalApi[site.SpiderKey]
	if !ok {
		easylog.Log.Error(source + " 没找到映射采集接口")
		return
	}
	// 调用采集函数
	res, err := f.ExternalGetById(site, VodId)
	if err != nil {
		easylog.Log.Error(err.Error())
		return
	}
	var datas []zyw.ZyMovie
	for _, data := range res.List {
		data.Source = source
		datas = append(datas, data)
	}
	saveDb(datas)
	cacheKey := source + VodId
	easycache.C.Delete(cacheKey)
}

// SearchHistory 最近搜索关键词
func SearchHistory(c *gin.Context) {
	history := viper.GetString("search_history")
	historys := strings.Split(history, ",")
	result.Ok(c, historys)
}

// DoubanImg 豆瓣图片代理
func DoubanImg(c *gin.Context) {
	link := c.Query("link")
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
	// 创建一个新的请求
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	// 添加 Header
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	// 发送请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	defer resp.Body.Close()
	// 设置响应头，告诉浏览器返回的是图片格式
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	// 将图片写入响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to write image"})
		return
	}
}

// 增加一条搜索记录
func addSearchHistory(kw string) {
	easylog.Log.Info("新增搜索记录：", kw)
	history := viper.GetString("search_history")
	historys := strings.Split(history, ",")
	if contains(historys, kw) {
		easylog.Log.Info("搜索记录已存在：", kw)
		return
	}
	historys = append(historys, kw)
	if len(historys) > 10 {
		// 只保留最新的10条记录
		historys = historys[len(historys)-10:]
	}
	str := strings.Join(historys, ",")
	easylog.Log.Info("最新的搜索记录：", str)
	// 修改配置文件
	viper.Set("search_history", str)
	if err := viper.WriteConfig(); err != nil {
		easylog.Log.Error(err.Error())
	}
}

// 判断字符串是否存在于字符串切片中
func contains(strs []string, s string) bool {
	for _, str := range strs {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}
