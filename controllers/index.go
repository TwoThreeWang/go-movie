package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/spider"
	"movie/utils/easylog"
	"movie/utils/try"
	"net/http"
)

func Index(c *gin.Context) {
	kw := c.Query("kw") // 搜索关键词
	result := gin.H{"data": "", "msg": "搜索词不能为空！", "status": http.StatusOK}
	if kw != "" {
		try.Try(func() {
			result["data"], result["msg"], result["status"] = spider.MacCmsSearch("https://pzdy.win", kw)
		}).CatchAll(func(err error) {
			easylog.Log.Error(err.Error())
			result["msg"] = err.Error()
			result["status"] = http.StatusBadRequest
		})
	}
	c.JSON(http.StatusOK, result)
}
