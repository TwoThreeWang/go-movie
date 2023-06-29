package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/spider"
	"movie/utils/auth"
	"movie/utils/easylog"
	"movie/utils/result"
	"movie/utils/try"
)

func Index(c *gin.Context) {
	// token校验
	if !auth.Verify(c) {
		c.Abort()
		return
	}
	kw := c.Query("kw") // 搜索关键词
	if kw != "" {
		try.Try(func() {
			data, msg := spider.MacCmsSearch("https://pzdy.win", kw)
			if msg != "" {
				result.Fail(c, result.ResultError, msg)
				c.Abort()
				return
			}
			result.Ok(c, data)
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
