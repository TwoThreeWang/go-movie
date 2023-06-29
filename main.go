package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	cache "github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"movie/routers"
	"movie/utils/easylog"
	"time"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

func cachedemo() {
	// 实例化缓存组件
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set("foo", "111111", cache.DefaultExpiration)
	if x, found := c.Get("foo"); found {
		foo := x.([]User)
		fmt.Println(foo)
	}
}

func dbdemo() {
	db, err := gorm.Open(sqlite.Open("data/sqlite.db?mode=wal"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}

	user := &User{Username: "alice", Password: "123456"}
	db.Create(&user)

	var users []User
	db.Find(&users)

	fmt.Println("users: ", users)
}

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
