package routers

import (
	"github.com/gin-gonic/gin"
	"movie/controllers/admin"
	"movie/utils/auth"
)

func RegisterAdminRoutes(r *gin.Engine) {
	api := r.Group("/admin")
	api.Use(auth.JWTAuth())
	{
		api.GET("/site", admin.ListSites)    // 列出所有站点
		api.POST("/site/add", admin.AddSite) // 新增站点
		api.GET("/site/del", admin.DelSite)  // 删除站点
	}
}
