package routers

import (
	"github.com/gin-gonic/gin"
	"movie/controllers"
)

func RegisterIndexRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		api.GET("/", controllers.Search)
		api.GET("/search", controllers.Search) // 搜索接口
	}
}
