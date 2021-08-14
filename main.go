package main

import (
	"fmt"
	"north-api/libs"
	"north-api/routers"
	"north-api/services"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var router *gin.Engine

func init() {
	LoadConfig()
	UseMiddlewares()
	InitRoutes()
	InitServices()
	InitLibs()
}

// LoadConfig 加载配置文件
func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading configuration file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

// InitRoutes 初始化路由
func InitRoutes() {
	router = gin.Default()
	// swagger
	routers.SetSwaggerRoutes(router)
	north := router.Group("/north")
	// 使用组中间件
	routers.UseNorthMiddlewares(north)
	routers.SetNorthRoutes(north)
}

// UseMiddlewares 使用全局中间件
func UseMiddlewares() {
}

// InitServices 初始化服务
func InitServices() {
	InitDb()
	InitRedis()
}

// InitDb 初始化数据库服务
func InitDb() {
	backend := viper.GetString("db.backend")
	dsn := viper.GetString("db.dsn")
	err := services.OpenDb(backend, dsn)
	if err != nil {
		panic(fmt.Errorf("Failed to open database: %s \n", err))
	}
	err = services.MigrateDb()
	if err != nil {
		panic(fmt.Errorf("Failed to migrate database: %s \n", err))
	}
}

// InitDb 初始化redis服务
func InitRedis() {
	var redisOptions redis.Options
	err := viper.UnmarshalKey("redis", &redisOptions)
	if err != nil {
		panic(fmt.Errorf("Read redis configuration error: %s \n", err))
	}
	services.OpenRedis(&redisOptions)
}

// InitLibs 初始化类库
func InitLibs() {
	libs.InitCache()
}

// Release 释放资源
func Release() {
	services.Db.Close()
}

// @title 北向API接口
// @version 1.0
// @description 北向API接口

// @license.name Apache 2.0

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	defer Release()
}
