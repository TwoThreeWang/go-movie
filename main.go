package main

import (
	"movie/routers"
	"movie/utils/easylog"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func configInit() bool {
	// 配置文件初始化
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		easylog.Log.Error("Error reading config file")
		return false
	}
	return true
}

func main() {
	// 日志初始化
	easylog.Init()
	easylog.Log.Info(`程序运行`)
	// 配置文件初始化
	if configInit() {
		// 将 Gin 的运行模式设置为 "release"，默认为debug模式
		gin.SetMode(gin.ReleaseMode)
		// 初始化Gin引擎
		r := gin.Default()
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		// 指定静态资源目录
		r.Static("/public", "./templates/public")
		r.StaticFile("/robots.txt", "public/robots.txt")
		//r.LoadHTMLGlob("templates/*/*.tmpl")
		routers.RegisterIndexRoutes(r)
		routers.RegisterAdminRoutes(r)
		routers.RegisterLoginRoutes(r)
		err := r.Run(":5005")
		if err != nil {
			easylog.Log.Info(err)
		}
	}
}
