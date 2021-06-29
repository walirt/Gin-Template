## 目录
```shell
.
├── config.yml # 配置文件
├── constant # 常量
│   ├── cache.go
│   ├── request.go
│   └── response.go
├── controllers # 控制器
│   └── north_c.go
├── docs # swagger文档相关
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── ent # orm相关
│   ... 
├── go.mod
├── go.sum
├── libs # 类库
│   └── cache.go
├── main.go # 主入口文件
├── middlewares # 中间件
│   └── north_m.go
├── README.md 
├── routers # 路由
│   ├── north_r.go
│   └── swagger_r.go
├── services # 服务
│   ├── db.go
│   └── redis.go
├── tmp # air临时文件
│   ├── build-errors.log
│   └── main
└── utils # 工具函数
    └── utils.go
```

## Web框架
[gin](https://github.com/gin-gonic/gin)

## ORM框架
[ent](https://github.com/ent/ent)
```shell
go run entgo.io/ent/cmd/ent init User # 初始化表(Schmea)
... # 编辑相关field, edges等
go generate ./ent # 生成表相关代码
```

## 配置文件读取
[viper](https://github.com/spf13/viper)

## 文档
[swaggo](https://github.com/swaggo/swag)
```shell
... # 编辑相关注释 
swag init # 生成相关代码
```

## 热重载
[air](https://github.com/cosmtrek/air)
```shell
air # 编译运行
```