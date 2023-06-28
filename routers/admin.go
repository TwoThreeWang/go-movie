package routers

import (
	"github.com/gin-gonic/gin"
	"movie/controllers/admin"
)

func RegisterAdminRoutes(r *gin.Engine) {
	api := r.Group("/admin")
	{
		api.GET("/", admin.ListSites)
		api.POST("/site/add", admin.AddSite)
	}
}
