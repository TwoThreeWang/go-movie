package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"movie/configs"
	"net/http"
)

func GetSites() (sites []configs.SpiderSites, err error) {
	// 获取已有站点数据
	err = viper.UnmarshalKey("sites", &sites)
	return
}

func AddSite(c *gin.Context) {
	result := gin.H{"data": "", "msg": "Success Add", "status": http.StatusOK}
	//绑定json参数和结构体
	var site configs.SpiderSites
	if err := c.BindJSON(&site); err != nil {
		result["msg"] = err.Error()
		result["status"] = http.StatusBadRequest
		c.AbortWithStatusJSON(http.StatusOK, result)
		return
	}
	result["data"] = site
	// 获取已有站点数据
	sites, err := GetSites()
	if err != nil {
		result["msg"] = err.Error()
		result["status"] = http.StatusBadRequest
		c.AbortWithStatusJSON(http.StatusOK, result)
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
		result["msg"] = err.Error()
		result["status"] = http.StatusBadRequest
		c.AbortWithStatusJSON(http.StatusOK, result)
		return
	}
	c.JSON(http.StatusOK, result)
}

func ListSites(c *gin.Context) {
	result := gin.H{"data": "", "msg": "Success", "status": http.StatusOK}
	// 获取已有站点数据
	sites, err := GetSites()
	result["data"] = sites
	if err != nil {
		result["msg"] = err.Error()
		result["status"] = http.StatusBadRequest
	}
	c.JSON(http.StatusOK, result)
}
