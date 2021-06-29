package controllers

import (
	"context"
	"fmt"
	"net/http"
	"north-api/constant"
	"north-api/ent/user"
	"north-api/libs"
	"north-api/services"
	"north-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/mssola/user_agent"
	"github.com/spf13/viper"
)

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept  json
// @Produce  json
// @Param reqestBody body constant.Request{version=string,data=object{username=string,password=string}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{login_time=string,timeout=int}} "登录成功"
// @Header 200 {string} Token "登录凭证"
// @Router /north/login [post]
func Login(c *gin.Context) {
	requestBody := c.MustGet("requestBody").(constant.Request)
	data := requestBody.Data.(map[string]interface{})
	usernameI, uok := data["username"]
	passwordI, pok := data["password"]
	if !(uok && pok) {
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.ParameterAbnormal,
			Data:          nil,
		})
		return
	}
	username := usernameI.(string)
	password := passwordI.(string)

	_, err := services.Db.User.
		Query().
		Where(user.UsernameEQ(username)).
		Where(user.PasswordEQ(utils.Md5Encrypt(password))).
		Only(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.UsernameOrPasswordWrong,
			Data:          nil,
		})
		return
	}

	allowMultipleSignIn := viper.GetBool("app.allowMultipleSignIn")
	// 不允许多点登录
	if !allowMultipleSignIn {
		tokens := make([]string, 0)
		// 清除登录列表信息
		listIter := libs.CacheScan(fmt.Sprintf(constant.LoggedInListKey, username, "*"), 0, 0)
		for listIter.Next(context.TODO()) {
			key := listIter.Val()
			var userInfo constant.UserInfo
			err := libs.CacheGet(key, &userInfo)
			if err != nil {
				c.JSON(http.StatusOK, constant.Response{
					ResponseError: constant.UnknownError,
					Data:          nil,
				})
				return
			}
			tokens = append(tokens, userInfo.Token)
			err = libs.CacheDelete(key)
			if err != nil {
				c.JSON(http.StatusOK, constant.Response{
					ResponseError: constant.UnknownError,
					Data:          nil,
				})
				return
			}
		}
		if err := listIter.Err(); err != nil {
			c.JSON(http.StatusOK, constant.Response{
				ResponseError: constant.UnknownError,
				Data:          nil,
			})
			return
		}
		// 清除登录状态
		for _, token := range tokens {
			err := libs.CacheDelete(fmt.Sprintf(constant.TokenKey, token))
			if err != nil {
				c.JSON(http.StatusOK, constant.Response{
					ResponseError: constant.UnknownError,
					Data:          nil,
				})
				return
			}
		}
	}

	// 保存登录状态
	token := utils.GenerateToken()
	tokenKey := fmt.Sprintf(constant.TokenKey, token)
	timeout := viper.GetInt64("app.tokenTimeOut")
	err = libs.CacheSet(tokenKey, username, time.Duration(timeout)*time.Second)
	if err != nil {
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.UnknownError,
			Data:          nil,
		})
		return
	}

	loggedInListKey := fmt.Sprintf(constant.LoggedInListKey, username, utils.RandString(6))
	ua := user_agent.New(c.GetHeader("User-Agent"))
	os := ua.OS()
	browerName, browerVersion := ua.Browser()
	// 保存登录列表信息
	userInfo := constant.UserInfo{
		Username:  username,
		Ip:        c.ClientIP(),
		Token:     token,
		Device:    fmt.Sprintf("%s, %s/%s", os, browerName, browerVersion),
		LoginTime: time.Now(),
	}
	err = libs.CacheSet(loggedInListKey, userInfo, time.Duration(timeout)*time.Second)
	if err != nil {
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.UnknownError,
			Data:          nil,
		})
		return
	}

	// 心跳选项
	heartbeatEnabled := viper.GetBool("app.heartbeatEnabled")
	if heartbeatEnabled {
		heartbeatKeepAliveTimeLimit := viper.GetInt("app.heartbeatKeepAliveTimeLimit")
		if heartbeatKeepAliveTimeLimit == 0 {
			heartbeatKeepAliveTimeLimit = 30
		}
		heartbeatKey := fmt.Sprintf(constant.HeartbeatKey, token)
		err := libs.CacheSet(heartbeatKey, true, time.Duration(heartbeatKeepAliveTimeLimit)*time.Second)
		if err != nil {
			c.JSON(http.StatusOK, constant.Response{
				ResponseError: constant.UnknownError,
				Data:          nil,
			})
			return
		}
	}

	c.Header("token", token)
	c.JSON(http.StatusOK, constant.Response{
		ResponseError: constant.OK,
		Data: map[string]interface{}{
			"login_time": time.Now().Format(time.RFC3339),
			"timeout":    timeout,
		},
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{logout_time=string}} "登出成功"
// @Router /north/logout [post]
func Logout(c *gin.Context) {
	token := c.GetHeader("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.TokenAbnormal,
			Data:          nil,
		})
		return
	}

	tokenKey := fmt.Sprintf(constant.TokenKey, token)
	var username string
	err := libs.CacheGet(tokenKey, &username)
	if err != nil && err != cache.ErrCacheMiss {
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.UnknownError,
			Data:          nil,
		})
		return
	}
	if err == cache.ErrCacheMiss {
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.TokenAbnormal,
			Data:          nil,
		})
		return
	}

	// 清除登录状态
	err = libs.CacheDelete(tokenKey)
	if err != nil {
		c.JSON(http.StatusOK, constant.Response{
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
				c.JSON(http.StatusOK, constant.Response{
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
			c.JSON(http.StatusOK, constant.Response{
				ResponseError: constant.UnknownError,
				Data:          nil,
			})
			return
		}
	}
	if err := listIter.Err(); err != nil {
		c.JSON(http.StatusOK, constant.Response{
			ResponseError: constant.UnknownError,
			Data:          nil,
		})
		return
	}
	c.JSON(http.StatusOK, constant.Response{
		ResponseError: constant.OK,
		Data: map[string]interface{}{
			"logout_time": time.Now().Format(time.RFC3339),
		},
	})
}

// Heartbeat 心跳
// @Summary 心跳
// @Description 心跳接口
// @Tags 通讯
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{heartbeat_time=string}}
// @Router /north/heartbeat [post]
func Heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, constant.Response{
		ResponseError: constant.OK,
		Data: map[string]interface{}{
			"heartbeat_time": time.Now().Unix(),
		},
	})
}

// ConfigGet 获取配置
// @Summary 获取配置
// @Description 获取配置接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Request{version=string,data=object{version=string,gzip=boolean}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{version=string,nodes=[]constant.SpaceObject}}
// @Router /north/config_get [post]
func ConfigGet(c *gin.Context) {
}

// OnlineDataGet 获取在线数据(请求/响应)
// @Summary 获取在线数据(请求/响应)
// @Description 获取在线数据(请求/响应)接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Request{version=string,data=object{space_guids=[]string,device_guids=[]string,point_guids=[]string}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{devices=[]constant.ResponseDevicesPart}}
// @Router /north/online_data_get [post]
func OnlineDataGet(c *gin.Context) {
}

// OnlineDataStrategyAdd 获取在线数据(发布/订阅) - 添加订阅策略
// @Summary 获取在线数据(发布/订阅) - 添加订阅策略
// @Description 获取在线数据(发布/订阅) - 添加订阅策略接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Request{version=string,data=object{name=string,space_guids=[]string,device_guids=[]string,point_guids=[]string,mode=int}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{strategy_id=int}}
// @Router /north/online_data_strategy_add [post]
func OnlineDataStrategyAdd(c *gin.Context) {
}

// OnlineDataStrategyDel 获取在线数据(发布/订阅) - 删除订阅策略
// @Summary 获取在线数据(发布/订阅) - 删除订阅策略
// @Description 获取在线数据(发布/订阅) - 删除订阅策略接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Request{version=string,data=object{strategy_id=int}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=string}
// @Router /north/online_data_strategy_del [post]
func OnlineDataStrategyDel(c *gin.Context) {
}

// OnlineDataStrategyQuery 获取在线数据(发布/订阅) - 查询订阅策略
// @Summary 获取在线数据(发布/订阅) - 查询订阅策略
// @Description 获取在线数据(发布/订阅) - 查询订阅策略接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Request{version=string,data=string} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{strategies=[]object{strategy_id=int,space_guids=[]string,device_guids=[]string,point_guids=[]string,mode=int}}}
// @Router /north/online_data_strategy_query [post]
func OnlineDataStrategyQuery(c *gin.Context) {
}

// OnlineDataStrategyPush 获取在线数据(发布/订阅) - 推送数据(客户端侧)
// @Summary 获取在线数据(发布/订阅) - 推送数据(客户端侧)
// @Description 获取在线数据(发布/订阅) - 推送数据(客户端侧)接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Response{version=string,data=object{strategy_id=int,devices=[]constant.ResponseDevicesPart}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=string}
// @Router /north/online_data_strategy_push [post]
func OnlineDataStrategyPush(c *gin.Context) {
}

// OfflineDataGet 获取离线数据
// @Summary 获取离线数据
// @Description 获取离线数据接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Request{version=string,data=object{strategy_id=int,begin_time=int,end_time=int}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{strategy_id=int,devices=[]constant.ResponseDevicesPart}}
// @Router /north/offline_data_get [post]
func OfflineDataGet(c *gin.Context) {
}

// HistoryDataGet 获取历史数据
// @Summary 获取历史数据
// @Description 获取历史数据接口
// @Tags 数据
// @Accept  json
// @Produce  json
// @Param token header string true "登录凭证"
// @Param reqestBody body constant.Request{version=string,data=object{point_guids=[]string,begin_time=int,end_time=int}} true "请求体"
// @Success 200 {object} constant.Response{error_code=int,error_msg=string,data=object{points=[]object{guid=string,tag=string,values=[]object{value=string,timestamp=int}}}}
// @Router /north/history_data_get [post]
func HistoryDataGet(c *gin.Context) {
}
