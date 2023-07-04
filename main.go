package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"movie/routers"
	"movie/utils/easylog"
)

func configInit() bool {
	// 配置文件初始化
	viper.SetConfigFile("configs/config.yaml")
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
		// 初始化Gin引擎
		r := gin.Default()
		routers.RegisterIndexRoutes(r)
		routers.RegisterAdminRoutes(r)
		routers.RegisterLoginRoutes(r)
		err := r.Run()
		if err != nil {
			easylog.Log.Info(err)
		}
	}
}
