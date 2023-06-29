package routers

import (
	"github.com/gin-gonic/gin"
	"movie/controllers"
)

func RegisterLoginRoutes(r *gin.Engine) {
	api := r.Group("/login")
	{
		api.POST("/", controllers.Login)
	}
}
