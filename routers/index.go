package routers

import (
	"github.com/gin-gonic/gin"
	"movie/controllers"
)

func RegisterIndexRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		api.GET("/", controllers.Index)
		api.GET("/search", controllers.Index) // 搜索接口
	}
}
