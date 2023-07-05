package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"movie/configs"
	"movie/repositories"
	"movie/services/spider"
	"movie/services/spider/douban"
	"movie/services/spider/zyw"
	"movie/utils/easycache"
	"movie/utils/easylog"
	"movie/utils/result"
	"movie/utils/try"
	"strconv"
	"sync"
	"time"
)

func GetSites() (sites []configs.SpiderSites, err error) {
	// 获取已有站点数据
	err = viper.UnmarshalKey("sites", &sites)
	return
}

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
						// 将响应结果发送到结果通道
						resultChan <- dataNew
						if err != nil {
							easylog.Log.Error(err.Error())
						}
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
						datas = append(datas, data)
					}
				case <-timeoutChan:
					cancel()
					fmt.Println("timeout!!!")
					if len(datas) > 0 {
						easycache.C.Set(kw, datas, cache.DefaultExpiration)
					}
					result.Ok(c, datas)
					c.Abort()
					return
				}
			}
			wg.Wait()
			if len(datas) > 0 {
				easycache.C.Set(kw, datas, cache.DefaultExpiration)
				go saveDb(datas) // 结果保存到数据库
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

func saveDb(datas []zyw.ZyMovie) {
	db := repositories.GetDB()
	var movies repositories.Movies
	for _, movie := range datas {
		movie.UpdatedAt = time.Now()
		movies = repositories.Movies(movie)
		res := db.Where("source = ? and vod_id = ?", movies.Source, movies.VodId).FirstOrCreate(&movies, movies)
		if res.Error != nil {
			easylog.Log.Error(res.Error)
		}
	}

}

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

func Index(c *gin.Context) {
	var data = map[string]string{
		"app_name": "api v1.0",
	}
	result.Ok(c, data)
}

func GetMovieInfo(c *gin.Context) {
	source := c.Query("source") // 来源
	VodId := c.Query("id")      // 影片ID
	db := repositories.GetDB()
	var movie repositories.Movies
	res := db.Where("source = ? and vod_id = ?", source, VodId).First(&movie)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			result.FailNoMsg(c, result.NotFound)
			c.Abort()
			return
		}
		result.FailNoMsg(c, result.DatabaseError)
		c.Abort()
		return
	}
	result.Ok(c, movie)
}
