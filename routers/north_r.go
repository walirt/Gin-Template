package routers

import (
	"north-api/controllers"
	"north-api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SetNorthRoutes 设置north路由 局部中间件
func SetNorthRoutes(router *gin.RouterGroup) {
	router.POST("login", controllers.Login)
	router.POST("logout", controllers.Logout)
	loginRequiredMiddleware := middlewares.LoginRequired()
	heartbeatDetectMiddleware := middlewares.HeartbeatDetect()
	partMiddlewaresList := make([]gin.HandlerFunc, 0)
	partMiddlewaresList = append(partMiddlewaresList, loginRequiredMiddleware)
	if viper.GetBool("app.heartbeatEnabled") {
		partMiddlewaresList = append(partMiddlewaresList, heartbeatDetectMiddleware)
	}
	{
		router.POST("heartbeat", append(partMiddlewaresList, controllers.Heartbeat)...)
		router.POST("config_get", append(partMiddlewaresList, controllers.ConfigGet)...)
	}

}

// UseNorthMiddlewares 设置north组中间件
func UseNorthMiddlewares(router *gin.RouterGroup) {
	router.Use(middlewares.IPBlocker())
	router.Use(middlewares.BindJsonRequestBody())
}
