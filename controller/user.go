package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/killsari/api/common"
	"github.com/killsari/api/model"
	"github.com/killsari/api/utils"
	"github.com/sirupsen/logrus"
)

type UserController struct {
}

func (UserController) Login(c *gin.Context) {
	ua := model.UserAgentModel{}.GetUA(c)
	uid := ua.Uid
	token := ua.Token

	type put struct {
		Tel  string `json:"tel" binding:"required,numeric,len=11"`
		Code string `json:"code" binding:"required,numeric,len=6"`
	}
	p := put{}
	if err := utils.ReqBindBody(c, &p); err != nil {
		return
	}
	// 校验验证码
	codeKey := fmt.Sprintf(common.RedisKeyCode, model.UserOptionRegLogin, p.Tel)
	code, err := common.Redis.Get(codeKey).Result()
	if err != nil {
		if err == redis.Nil {
			utils.ResErrWithMsg(c, "验证码不存在或已过期")
			return
		}
		logrus.Error(err)
		utils.ResErrDefault(c)
		return
	}
	if code != p.Code {
		utils.ResErrWithMsg(c, "验证码错误")
		return
	}
	// 用户
	userDao, err := model.UserModel{}.GetByTel(p.Tel)
	newFlag := false
	if err != nil {
		if err != sql.ErrNoRows {
			utils.ResErrDefault(c)
			return
		}
		// 新用户注册
		newFlag = true
		userDao = &model.UserDao{
			Tel:    p.Tel,
		}
		err := model.UserModel{}.Put(userDao)
		if err != nil {
			utils.ResErrDefault(c)
			return
		}
	}
	// 验证后删除缓存
	if err = common.Redis.Del(codeKey).Err(); err != nil {
		logrus.Error(err) // ignore error output
	}
	// 登陆设置token
	if userDao.Id != uid { //已登陆用户再次登陆自己账户，不生成新token
		token = utils.TokenGen(userDao.Tel)
		tokenKey := fmt.Sprintf(common.RedisKeyToken, token)
		err := common.Redis.Set(tokenKey, userDao.Id, common.RedisTtlToken).Err()
		if err != nil {
			logrus.Error(err)
			utils.ResErrDefault(c)
			return
		}
	}
	utils.SetAuthorization(c, token)
	// 用户详情,新用户需要重新获取用户信息，不可用Show()
	userShow := &model.UserShow{}
	if newFlag {
		userShow, err = model.UserModel{}.ShowById(userDao.Id, userDao.Id, false)
	} else {
		userShow, err = model.UserModel{}.Show(userDao, userDao.Id, false)
	}
	if err != nil {
		utils.ResErrDefault(c)
		return
	}

	utils.ResSuc(c, userShow)
}

func (UserController) Logout(c *gin.Context) {
	token := model.UserAgentModel{}.GetToken(c)
	tokenKey := fmt.Sprintf(common.RedisKeyToken, token)
	err := common.Redis.Del(tokenKey).Err()
	if err != nil {
		logrus.Error(err)
		utils.ResErrDefault(c)
		return
	}

	utils.ResSucWithoutData(c)
}

func (UserController) Get(c *gin.Context) {
	uid := model.UserAgentModel{}.GetUid(c)
	show, err := model.UserModel{}.ShowById(uid, uid, false)
	if err != nil {
		utils.ResErrDefault(c)
		return
	}
	utils.ResSuc(c, show)
}
