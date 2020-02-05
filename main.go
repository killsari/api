package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/killsari/api/common"
	"github.com/killsari/api/controller"
	"github.com/killsari/api/model"
	_ "github.com/killsari/api/service"
	"github.com/killsari/api/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	r := setupRouter()
	r.Run(":8083")
}

func setupRouter() *gin.Engine {
	r := gin.Default() //== gin.New().Use(Logger(), Recovery())
	//中间件|过滤器|...:
	r.Use(cores())             //跨域控制
	r.Use(regularMiddleware()) //常规认证
	r.Use(authMiddleware())    //用户认证
	r.Use(uaMiddleware())      //用户设备
	//接口:
	pub := r.Group("/pub") //公共
	{
		pub.GET("/config", controller.PubController{}.Config)  //一些配置信息
		pub.POST("/code", controller.PubController{}.Code)     //注册/登陆验证码
		pub.GET("/region", controller.PubController{}.Region)  //行政区划x3
	}
	user := r.Group("/user") //用户
	{
		user.POST("/login", controller.UserController{}.Login)   //登陆/注册
		user.GET("/get", controller.UserController{}.Get)        //获取用户信息
		user.POST("/logout", controller.UserController{}.Logout) //退登
	}

	return r
}

/**
授权过滤器
*/
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取授权用户
		pass, err := getTokenUserAndCheck(c)
		if err != nil {
			utils.ResErr(c, common.ErrCodeSys, err.Error())
			c.Abort()
			return
		}
		if pass {
			c.Next()
			return
		}
		//过滤需要授权的url
		noAuthUriArr := []string{
			"/pub/config", "/pub/code", "/pub/region",
			"/user/login",
		}
		path := strings.Split(c.Request.RequestURI, "?")[0]
		for _, item := range noAuthUriArr {
			if path == item {
				c.Next()
				return
			}
		}
		c.Abort()
		utils.ResErr(c, common.ErrCodeAuth, "auth error")
	}
}

/**
获取授权用户并且合法性验证
*/
func getTokenUserAndCheck(c *gin.Context) (pass bool, err error) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "" {
		arr := strings.Split(authorization, " ")
		if len(arr) == 2 && arr[0] == "Bearer" {
			token := arr[1]
			redisTokenKey := fmt.Sprintf(common.RedisKeyToken, token)
			var uidCache string
			uidCache, err = common.Redis.Get(redisTokenKey).Result()
			if err == nil {
				//extend token when needed
				var duration time.Duration
				duration, err = common.Redis.TTL(redisTokenKey).Result()
				if err == nil {
					if duration < common.RedisTtlExtendToken {
						if err = common.Redis.Set(redisTokenKey, uidCache, common.RedisTtlToken).Err(); err != nil {
							logrus.Error(err)
							return
						}
					}
					var uid uint64
					uid, err = strconv.ParseUint(uidCache, 10, 64)
					if err != nil {
						logrus.Error(err)
						return
					}
					// -- 所有token代表有效用户，后台锁定和删除用户会清除token
					var ua model.UserAgentEntity
					value, exists := c.Get("ua")
					if exists {
						ua = value.(model.UserAgentEntity)
					}
					ua.Uid = uid
					ua.Token = token
					c.Set("ua", &ua)
					pass = true
				} else if err == redis.Nil {
					err = nil //do not return fail msg
				} else {
					logrus.Error(err)
					return
				}
			} else if err == redis.Nil {
				err = nil //do not return fail msg
			} else {
				logrus.Error(err)
				return
			}
		}
	}
	return
}

/**
授权过滤器
*/
func uaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uaStr := c.Request.Header.Get(model.UserAgentHeader)
		if uaStr != "" {
			value, exists := c.Get("ua")
			var ua model.UserAgentEntity
			if exists {
				uaP := value.(*model.UserAgentEntity)
				ua = *uaP
			}
			uaArr := strings.Split(uaStr, model.UserAgentSep)
			ua.Language = uaArr[0]
			if len(uaArr) > 1 {
				ua.Source = uaArr[1]
			}
			if len(uaArr) > 2 {
				ua.Version = uaArr[2]
			}
			if len(uaArr) > 3 {
				ua.Resolution = uaArr[3]
			}
			if len(uaArr) > 4 {
				ua.Device = uaArr[4]
			}
			if len(uaArr) > 5 {
				ua.Device = uaArr[5]
			}
			c.Set("ua", &ua)
		}
		c.Next()
		return
	}
}

/**
公共中间件
*/
func regularMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//接口状态全局处理
		status := c.Writer.Status()
		if status != http.StatusOK {
			utils.ResErr(c, uint64(status), "http error")
			c.Abort()
			return
		}
		//接口响应时间处理start
		t := time.Now()
		//before req
		c.Next()
		//after req
		//接口响应时间处理end
		latency := time.Since(t)
		if latency > time.Millisecond*200 {
			logrus.Warn(fmt.Sprintf("%s cost %s", c.Request.URL, latency.String()))
		}
	}
}

func cores() gin.HandlerFunc {
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*", "Content-Type", "Authorization", model.UserAgentHeader},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AddExposeHeaders("*", "Content-Type", "Authorization", model.UserAgentHeader)
	handlerFunc := cors.New(config)
	return handlerFunc
}
