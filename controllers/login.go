package controllers

import (
	"github.com/gin-gonic/gin"
	"movie/middlewares"
	"movie/utils/result"
)

func Login(c *gin.Context) {
	var user middlewares.LoginUser
	if err := c.BindJSON(&user); err != nil {
		result.FailNoMsg(c, result.InvalidArgs)
		c.Abort()
		return
	}
	// todo 数据库查询用户名密码是否正确
	signedToken, err := middlewares.GetToken(user)
	if err != nil {
		result.Fail(c, result.ResultError, err.Error())
		c.Abort()
		return
	}
	data := map[string]any{
		"token":       signedToken,
		"expire_time": middlewares.ExpireTime,
	}
	result.Ok(c, data)
}
