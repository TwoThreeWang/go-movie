package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/spider"
	"movie/utils/easylog"
	"movie/utils/result"
	"movie/utils/try"
)

func Index(c *gin.Context) {
	kw := c.Query("kw") // 搜索关键词
	if kw != "" {
		try.Try(func() {
			data, msg := spider.MacCmsSearch("https://pzdy.win", kw)
			if msg != "" {
				result.Fail(c, result.ResultError, msg)
				return
			}
			result.Ok(c, data)
			return
		}).CatchAll(func(err error) {
			easylog.Log.Error(err.Error())
			result.Fail(c, result.ResultError, err.Error())
			return
		})
	}
	result.FailNoMsg(c, result.InvalidArgs)
}
