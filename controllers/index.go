package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"movie/configs"
	"movie/spider"
	"movie/spider/maccms"
	"movie/utils/auth"
	"movie/utils/easylog"
	"movie/utils/result"
	"movie/utils/try"
	"reflect"
	"strconv"
	"time"
)

// 定义一个map，存储不同key对应的struct类型
var types = map[string]reflect.Type{
	"maccms": reflect.TypeOf(maccms.MacCmsSearch{}), // 苹果CMS采集接口
}

func GetSites() (sites []configs.SpiderSites, err error) {
	// 获取已有站点数据
	err = viper.UnmarshalKey("sites", &sites)
	return
}

func Search(c *gin.Context) {
	// token校验
	if !auth.Verify(c) {
		c.Abort()
		return
	}
	kw := c.Query("kw") // 搜索关键词
	if kw != "" {
		try.Try(func() {
			sites, err := GetSites()
			if err != nil {
				result.Fail(c, result.ResultError, err.Error())
				c.Abort()
				return
			}
			// 创建结果通道
			resultChan := make(chan []spider.Movie, len(sites))
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
			for _, site := range sites {
				go func(ctx context.Context) {
					// 动态调用采集函数
					dataNew, msg := DynamicRunFunc(site, kw)
					// 将响应结果发送到结果通道
					resultChan <- dataNew
					if msg != "" {
						easylog.Log.Error(msg)
					}
				}(ctx)
			}
			// 等待所有响应结果返回或者超时
			var datas []spider.Movie
			for i := 0; i < len(sites); i++ {
				select {
				case res := <-resultChan:
					datas = append(datas, res...)
				case <-timeoutChan:
					fmt.Println("timeout!!!")
					result.Ok(c, datas)
					c.Abort()
					return
				}
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
	}
}

// DynamicRunFunc 动态调用不同的采集接口
func DynamicRunFunc(site configs.SpiderSites, kw string) (data []spider.Movie, msg string) {
	key := site.SpiderKey
	// 根据key获取对应的struct类型
	ts, ok := types[key]
	if !ok {
		msg = "unknown type " + key
		return
	}
	// 动态实例化struct
	sp := reflect.New(ts).Elem()
	// 调用接口方法
	spd, ok := sp.Interface().(spider.Spider)
	if !ok {
		msg = key + " does not implement Shape interface"
		return
	}
	// 调用接口中的Search函数
	return spd.SpiderSearch(site, kw)
}
