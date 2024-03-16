# go-movie
go-movie

使用 Gin 框架编写一个完整的 API 项目，目录结构参考：
```
├── cmd
│   └── myapp
│       └── main.go
├── configs
│   ├── config.go
│   └── database.go
├── controllers
│   ├── auth_controller.go
│   └── user_controller.go
├── middlewares
│   ├── auth_middleware.go
│   └── logger_middleware.go
├── models
│   ├── user.go
│   └── token.go
├── repositories
│   ├── user_repository.go
│   └── token_repository.go
├── routes
│   ├── auth_routes.go
│   └── user_routes.go
├── services
│   ├── auth_service.go
│   └── user_service.go
└── go.mod
```
下面是每个目录的作用：

cmd: 存放应用程序的启动代码，例如 main.go。

configs: 存放应用程序的配置信息，例如数据库连接信息等。

controllers: 存放控制器代码，处理 HTTP 请求并返回响应。

middlewares: 存放中间件代码，例如身份验证、日志记录等。

models: 存放数据模型的定义，例如用户、令牌等。

repositories: 存放数据访问对象的实现，与数据库进行交互。

routes: 存放路由定义和路由处理函数。

services: 存放服务层代码，封装业务逻辑并调用数据访问对象进行数据操作。

在这个目录结构中，控制器层负责处理 HTTP 请求和响应，服务层负责封装业务逻辑和调用数据访问对象进行数据操作，数据访问对象负责与数据库进行交互，中间件层负责处理请求前后的逻辑，例如身份验证、日志记录等。这种分层架构可以使代码更加清晰、易于维护和扩展。

mac下将GO打包成linux可执行文件命令：

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go_movie
```