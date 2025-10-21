package main

import (
	"movie/routers"
	"movie/utils/easylog"
	"net/http"
	"time"

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
		r.Use(gin.Recovery()) // 捕获 panic，防止服务崩溃
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		// 指定静态资源目录
		r.Static("/public", "./templates/public")
		r.StaticFile("/robots.txt", "public/robots.txt")
		//r.LoadHTMLGlob("templates/*/*.tmpl")
		routers.RegisterIndexRoutes(r)
		routers.RegisterAdminRoutes(r)
		routers.RegisterLoginRoutes(r)
		srv := &http.Server{
			Addr:              ":5005",
			Handler:           r,
			ReadTimeout:       30 * time.Second,  // 限制读取完整请求的时间（包括Body）
			WriteTimeout:      30 * time.Second,  // 限制写入响应的时间
			IdleTimeout:       120 * time.Second, // 限制空闲连接的保持时间（Keep-Alive）
			ReadHeaderTimeout: 10 * time.Second,  // 限制读取请求头的时间
			MaxHeaderBytes:    1 << 20,           // 1MB，限制请求头的最大大小
		}

		// 启动服务
		easylog.Log.Info("启动http服务,端口:5005,监听请求中...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}
}
