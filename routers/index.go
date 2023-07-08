package routers

import (
	"github.com/gin-contrib/multitemplate"
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
	r.HTMLRender = createMyRender()
	web := r.Group("/")
	{ // 首页
		web.GET("/", func(c *gin.Context) {
			c.HTML(200, "home", nil)
		})
		// 搜索结果页
		web.GET("/search", func(c *gin.Context) {
			c.HTML(200, "search", gin.H{
				"kw": c.Query("kw"),
			})
		})
		// 影片详情页
		web.GET("/detail", func(c *gin.Context) {
			c.HTML(200, "detail", gin.H{
				"source": c.Query("source"),
				"vid":    c.Query("vid"),
			})
		})
		// 播放页
		web.GET("/play", func(c *gin.Context) {
			c.HTML(200, "play", gin.H{
				"source": c.Query("source"),
				"vid":    c.Query("vid"),
				"play":   c.Query("play"),
			})
		})
	}
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("home", "templates/web/base.tmpl", "templates/web/header.tmpl",
		"templates/web/footer.tmpl", "templates/web/home.tmpl")
	r.AddFromFiles("search", "templates/web/base.tmpl", "templates/web/header.tmpl",
		"templates/web/footer.tmpl", "templates/web/search.tmpl")
	r.AddFromFiles("detail", "templates/web/base.tmpl", "templates/web/header.tmpl",
		"templates/web/footer.tmpl", "templates/web/detail.tmpl")
	r.AddFromFiles("play", "templates/web/base.tmpl", "templates/web/header.tmpl",
		"templates/web/footer.tmpl", "templates/web/play.tmpl")
	return r
}
