package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"movie/configs"
	"movie/utils/result"
	"strconv"
)

func GetSites() (sites []configs.SpiderSites, err error) {
	// 获取已有站点数据
	err = viper.UnmarshalKey("sites", &sites)
	return
}

// AddSite 新增站点
func AddSite(c *gin.Context) {
	//绑定json参数和结构体
	var site configs.SpiderSites
	if err := c.BindJSON(&site); err != nil {
		result.Fail(c, result.ResultError, err.Error())
		c.Abort()
		return
	}
	// 获取已有站点数据
	sites, err := GetSites()
	if err != nil {
		result.Fail(c, result.ResultError, err.Error())
		c.Abort()
		return
	}
	// 保存新增站点数据
	sites = append(sites, configs.SpiderSites{
		Key:       site.Key,
		BaseUrl:   site.BaseUrl,
		PlayUrl:   site.PlayUrl,
		SpiderKey: site.SpiderKey,
	})
	// 修改配置文件
	viper.Set("sites", sites)
	if err = viper.WriteConfig(); err != nil {
		result.Fail(c, result.ResultError, err.Error())
		c.Abort()
		return
	}
	result.Ok(c, site)
}

func ListSites(c *gin.Context) {
	// 获取已有站点数据
	sites, err := GetSites()
	if err != nil {
		result.Fail(c, result.ResultError, err.Error())
		c.Abort()
		return
	}
	result.Ok(c, sites)
}

// DelSite 删除站点
func DelSite(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		result.FailNoMsg(c, result.InvalidArgs)
		c.Abort()
		return
	}
	// 获取已有站点数据
	sites, err := GetSites()
	if err != nil {
		result.Fail(c, result.ResultError, err.Error())
		c.Abort()
		return
	}
	// 删除新增站点数据
	var newsites []configs.SpiderSites
	delnum := 0
	for _, site := range sites {
		if site.Key != key {
			newsites = append(newsites, site)
		} else {
			delnum += 1
		}
	}
	// 修改配置文件
	viper.Set("sites", newsites)
	if err = viper.WriteConfig(); err != nil {
		result.Fail(c, result.DatabaseError, err.Error())
		c.Abort()
		return
	}
	result.Ok(c, "成功删除 "+strconv.Itoa(delnum)+" 条数据")
}
