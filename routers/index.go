package routers

import (
	"github.com/gin-gonic/gin"
	"movie/controllers"
)

func RegisterIndexRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/", controllers.Index)
		api.GET("/search", controllers.Search)               // 搜索接口
		api.GET("/douban", controllers.DoubanHot)            // 豆瓣热搜
		api.GET("/detail", controllers.GetMovieInfo)         // 影片详情
		api.GET("/searchhistory", controllers.SearchHistory) // 最近搜索关键词
	}
	web := r.Group("/")
	{
		web.GET("/", func(c *gin.Context) {
			c.HTML(200, "home.html", gin.H{
				"s": "这是首页，home",
			})
		})
	}
}
