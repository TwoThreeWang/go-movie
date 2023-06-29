package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一返回报文结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	OK              = Error(200, "success")              // 成功
	NeedRedirect    = Error(301, "need redirect")        // 重定向
	InvalidArgs     = Error(400, "invalid params")       // 参数错误
	Unauthorized    = Error(401, "unauthorized")         // 未授权
	Forbidden       = Error(403, "forbidden")            // 禁止
	NotFound        = Error(404, "not found")            // 404未找到
	Conflict        = Error(409, "entry exist")          // 冲突，条目存在
	TooManyRequests = Error(429, "too many requests")    // 请求过多
	ResultError     = Error(500, "response error")       // 响应错误
	DatabaseError   = Error(598, "database error")       // 数据错误
	CSRFDetected    = Error(599, "csrf attack detected") // CSRF 攻击

	UserError  = Error(5001, "username or password error") // 用户名或密码错误
	CodeExpire = Error(5002, "verification expire")        // 验证到期
	CodeError  = Error(5003, "verification error")         // 验证错误
	UserExist  = Error(5004, "user Exist")                 // 用户退出
)

func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}
func Error(code int, msg string) Response {
	return Response{
		code,
		msg,
		nil,
	}
}

func OkNoData(c *gin.Context) {
	Ok(c, nil)
}

func Ok(c *gin.Context, data interface{}) {
	Result(c, OK.Code, OK.Msg, data)
}

func OkWithMsg(c *gin.Context, msg string) {
	Result(c, OK.Code, msg, nil)
}

func FailNoMsg(c *gin.Context, err Response) {
	Result(c, err.Code, err.Msg, nil)
}

func Fail(c *gin.Context, err Response, msg string) {
	Result(c, err.Code, msg, nil)
}

//错误响应调用方式
//r.GET("/hello", func(c *gin.Context) {
//	result.Fail(c, result.ResultError)
//})
