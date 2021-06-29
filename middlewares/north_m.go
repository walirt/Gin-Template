package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"north-api/constant"
	"north-api/libs"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/spf13/viper"
)

func IPBlocker() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipWhitelist := viper.GetStringSlice("app.ipWhitelist")
		clientIp := c.ClientIP()
		for _, ip := range ipWhitelist {
			ok, _ := filepath.Match(ip, clientIp)
			if ok {
				c.Next()
				return
			}
		}
		responseError := constant.OtherError
		responseError.ErrorMsg = fmt.Sprintf(responseError.ErrorMsg, "ip is not in the whitelist")
		c.AbortWithStatusJSON(http.StatusOK, constant.Response{
			ResponseError: responseError,
			Data:          nil,
		})
	}
}

func BindJsonRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody constant.Request
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, constant.Response{
				ResponseError: constant.ParameterAbnormal,
				Data:          nil,
			})
			return
		}
		c.Set("requestBody", requestBody)
		c.Next()
	}
}

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, constant.Response{
				ResponseError: constant.TokenAbnormal,
				Data:          nil,
			})
			return
		}

		tokenKey := fmt.Sprintf(constant.TokenKey, token)
		var username string
		err := libs.CacheGet(tokenKey, &username)
		if err != nil && err != cache.ErrCacheMiss {
			c.AbortWithStatusJSON(http.StatusOK, constant.Response{
				ResponseError: constant.UnknownError,
				Data:          nil,
			})
			return
		}
		if err == cache.ErrCacheMiss {
			c.AbortWithStatusJSON(http.StatusOK, constant.Response{
				ResponseError: constant.TokenAbnormal,
				Data:          nil,
			})
			return
		}

		c.Set("currentUser", username)
		c.Next()
	}
}

func HeartbeatDetect() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, constant.Response{
				ResponseError: constant.TokenAbnormal,
				Data:          nil,
			})
			return
		}
		heartbeatKey := fmt.Sprintf(constant.HeartbeatKey, token)
		if libs.CacheExists(heartbeatKey) {
			// 在保活时限内
			heartbeatKeepAliveTimeLimit := viper.GetInt("app.heartbeatKeepAliveTimeLimit")
			if heartbeatKeepAliveTimeLimit == 0 {
				heartbeatKeepAliveTimeLimit = 30
			}
			err := libs.CacheSet(heartbeatKey, true, time.Duration(heartbeatKeepAliveTimeLimit)*time.Second)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, constant.Response{
					ResponseError: constant.UnknownError,
					Data:          nil,
				})
				return
			}
			c.Next()
			return
		} else {
			// 超过保活时限
			tokenKey := fmt.Sprintf(constant.TokenKey, token)
			var username string
			err := libs.CacheGet(tokenKey, &username)
			if err != nil && err != cache.ErrCacheMiss {
				c.AbortWithStatusJSON(http.StatusOK, constant.Response{
					ResponseError: constant.UnknownError,
					Data:          nil,
				})
				return
			}
			if err == cache.ErrCacheMiss {
				c.AbortWithStatusJSON(http.StatusOK, constant.Response{
					ResponseError: constant.TokenAbnormal,
					Data:          nil,
				})
				return
			}

			// 清除登录状态
			err = libs.CacheDelete(tokenKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, constant.Response{
					ResponseError: constant.UnknownError,
					Data:          nil,
				})
				return
			}

			allowMultipleSignIn := viper.GetBool("app.allowMultipleSignIn")
			// 清除登录列表信息
			removeLoggedInListKeys := make([]string, 0)
			listIter := libs.CacheScan(fmt.Sprintf(constant.LoggedInListKey, username, "*"), 0, 0)
			for listIter.Next(context.TODO()) {
				key := listIter.Val()
				if allowMultipleSignIn {
					var userInfo constant.UserInfo
					err := libs.CacheGet(key, &userInfo)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusOK, constant.Response{
							ResponseError: constant.UnknownError,
							Data:          nil,
						})
						return
					}
					if userInfo.Token == token {
						removeLoggedInListKeys = append(removeLoggedInListKeys, key)
						break
					}
				} else {
					removeLoggedInListKeys = append(removeLoggedInListKeys, key)
				}
			}
			for _, removeKey := range removeLoggedInListKeys {
				err = libs.CacheDelete(removeKey)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusOK, constant.Response{
						ResponseError: constant.UnknownError,
						Data:          nil,
					})
					return
				}
			}
			if err := listIter.Err(); err != nil {
				c.AbortWithStatusJSON(http.StatusOK, constant.Response{
					ResponseError: constant.UnknownError,
					Data:          nil,
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusOK, constant.Response{
				ResponseError: constant.TokenAbnormal,
				Data:          nil,
			})
			return
		}
	}
}
