package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"movie/utils/result"
	"time"
)

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
	Username string `json:"username"`
}

var (
	Secret     = "wangtwothree" // 加盐
	ExpireTime = 3600 * 12      // token有效期
)

const (
	ErrorReason_ServerBusy = "服务器繁忙"
	ErrorReason_ReLogin    = "请重新登陆"
)

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Verify(c *gin.Context) bool {
	// 通过http header中的token解析来认证
	strToken := c.Request.Header.Get("token")
	if strToken == "" {
		result.FailNoMsg(c, result.CodeError)
		c.Abort()
		return false
	}
	claim, err := verifyAction(strToken)
	if err != nil {
		result.Fail(c, result.Unauthorized, err.Error())
		c.Abort()
		return false
	}
	//easylog.Log.Info("登录用户：" + claim.Username)
	// 将解析后的有效载荷claims重新写入gin.Context引用对象中
	c.Set("claims", claim)
	return true
}

func verifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, errors.New(ErrorReason_ServerBusy)
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReason_ReLogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorReason_ReLogin)
	}
	return claims, nil
}

func GetToken(user LoginUser) (string, error) {
	// 构造用户信息
	claims := &JWTClaims{
		UserID:   1,
		Username: user.Username,
		Password: user.Password,
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	// 生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// JWTAuth 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token校验
		if !Verify(c) {
			c.Abort()
			return
		}
	}
}
