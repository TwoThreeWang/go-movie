package routers

import (
	"movie/controllers"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
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
				"source":  c.Query("source"),
				"vid":     c.Query("vid"),
				"refresh": c.Query("refresh"),
			})
		})
		// 播放页
		web.GET("/play", func(c *gin.Context) {
			c.HTML(200, "play", gin.H{
				"source": c.Query("source"),
				"vid":    c.Query("vid"),
				"play":   c.Query("play"),
				"name":   c.Query("name"),
			})
		})
		// 关于页
		web.GET("/about", func(c *gin.Context) {
			c.HTML(200, "about", gin.H{})
		})
		web.GET("/hot/:name", func(c *gin.Context) {
			c.HTML(200, "hot", gin.H{
				"name": c.Param("name"),
			})
		})
		web.GET("/page/:name", controllers.RenderMarkdownFile) // 文档解析
		// 豆瓣图片代理
		web.GET("/doubanimg", controllers.DoubanImg)
	}
	// 加载404错误页面
	r.NoRoute(func(c *gin.Context) {
		// 实现内部重定向
		c.HTML(200, "404", gin.H{
			"title": "404 - 页面未找到",
			"desc":  "抱歉，您要找的页面可能已经离场了。不如去看看其他精彩内容？",
		})
	})
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("home", "templates/web/base.html", "templates/web/home.html")
	r.AddFromFiles("search", "templates/web/base.html", "templates/web/search.html")
	r.AddFromFiles("detail", "templates/web/base.html", "templates/web/detail.html")
	r.AddFromFiles("play", "templates/web/base.html", "templates/web/play.html")
	r.AddFromFiles("about", "templates/web/base.html", "templates/web/about.html")
	r.AddFromFiles("page", "templates/web/base.html", "templates/web/page.html")
	r.AddFromFiles("hot", "templates/web/base.html", "templates/web/hot.html")
	r.AddFromFiles("404", "templates/web/base.html", "templates/web/result.html")
	return r
}
