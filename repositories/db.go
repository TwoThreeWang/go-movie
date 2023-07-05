package repositories

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

var db *gorm.DB

func init() {
	var err error
	// mysql 数据库连接，测试环境临时用 sqlite 代替
	//config := DatabaseConfig{
	//	Type:     viper.GetString("mysql.type"),
	//	Host:     viper.GetString("mysql.host"),
	//	Port:     viper.GetInt("mysql.port"),
	//	User:     viper.GetString("mysql.user"),
	//	Password: viper.GetString("mysql.password"),
	//	Name:     viper.GetString("mysql.dbname"),
	//}
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//	config.User, config.Password, config.Host, config.Port, config.Name)
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", "data/sqlite.db") // sqlite
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)  // 连接池最小连接数
	sqlDB.SetMaxOpenConns(100) // 连接池最大连接数
	// 自动迁移
	if err := db.AutoMigrate(&Movies{}); err != nil {
		fmt.Println(err.Error())
	}
}

func GetDB() *gorm.DB {
	return db
}
