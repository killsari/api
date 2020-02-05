package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/killsari/api/common"
	"github.com/killsari/api/model"
	"github.com/killsari/api/utils"
	"github.com/sirupsen/logrus"
)

type PubController struct {
}

func (controller PubController) Code(c *gin.Context) {
	type put struct {
		Tel string `json:"tel" binding:"required,numeric,len=11"`
	}
	p := put{}
	if err := utils.ReqBindBody(c, &p); err != nil {
		return
	}

	codeKey := fmt.Sprintf(common.RedisKeyCode, model.UserOptionRegLogin, p.Tel)
	err := common.Redis.Get(codeKey).Err()
	if err != nil {
		if err != redis.Nil {
			utils.ResErrDefault(c)
			return
		}
	} else {
		utils.ResErrWithMsg(c, "发送验证码过于频繁，请稍后再试")
		return
	}

	code := utils.CodeGen()

	//todo 上线需要释放以下发短信代码
	//err = utils.AliyunMsgUtils{}.SendSms(p.Tel, code)
	//if err != nil {
	//	utils.ResErrWithMsg(c, "发送验证码失败，请稍后重试")
	//	return
	//}

	err = common.Redis.Set(codeKey, code, common.RedisTtlCode).Err()
	if err != nil {
		utils.ResErrDefault(c)
		return
	}

	utils.ResSuc(c, map[string]string{"debug": code}) // todo 上线需要删除debug消息
}

func (controller PubController) Config(c *gin.Context) {
	result := map[string]string{
		"webUrl":      fmt.Sprintf("%s", common.WebUrlPrefix),
		"downloadUrl": fmt.Sprintf("%s/download", common.WebUrlPrefix),
	}
	utils.ResSuc(c, result)
}

func (controller PubController) Region(c *gin.Context) {
	type itemStruct struct {
		Code uint64                 `json:"code"`
		Name string                 `json:"name"`
		Cld  map[uint64]interface{} `json:"cld"`
	}
	regionMap := map[uint64]interface{}{}

	redisValue, err := common.Redis.Get(common.RedisKeyRegion).Result()
	if err != nil {
		if err != redis.Nil {
			utils.ResErrDefault(c)
			return
		}
		daoArr, err := model.RegionModel{}.Select()
		if err != nil {
			utils.ResErrDefault(c)
			return
		}

		//组合map
		for _, item := range *daoArr {
			provinceCodePre := item.Code / 10000
			provinceCode := provinceCodePre * 10000
			cityCodePre := item.Code % 10000 / 100
			cityCode := cityCodePre*100 + provinceCode
			districtCodePre := item.Code % 100
			districtCode := item.Code%100 + cityCode
			var ok bool
			if cityCodePre == 0 && districtCodePre == 0 {
				provinceCld := map[uint64]interface{}{}
				_, ok = regionMap[provinceCode]
				if ok {
					provinceCld = regionMap[provinceCode].(itemStruct).Cld
				}
				regionMap[provinceCode] = itemStruct{
					Code: item.Code,
					Name: item.Name,
					Cld:  provinceCld,
				}
			} else if districtCodePre == 0 {
				_, ok = regionMap[provinceCode]
				if !ok {
					regionMap[provinceCode] = itemStruct{
						Code: provinceCode,
						Name: "",
						Cld:  map[uint64]interface{}{},
					}
				}
				cityCld := map[uint64]interface{}{}
				_, ok = regionMap[provinceCode].(itemStruct).Cld[cityCode]
				if ok {
					cityCld = regionMap[provinceCode].(itemStruct).Cld[cityCode].(itemStruct).Cld
				}
				regionMap[provinceCode].(itemStruct).Cld[cityCode] = itemStruct{
					Code: item.Code,
					Name: item.Name,
					Cld:  cityCld,
				}
			} else {
				_, ok = regionMap[provinceCode]
				if !ok {
					regionMap[provinceCode] = itemStruct{
						Code: 0,
						Name: "",
						Cld:  map[uint64]interface{}{},
					}
				}
				_, ok = regionMap[provinceCode].(itemStruct).Cld[cityCode]
				if !ok {
					regionMap[provinceCode].(itemStruct).Cld[cityCode] = itemStruct{
						Code: cityCode,
						Name: "",
						Cld:  map[uint64]interface{}{},
					}
				}
				_, ok = regionMap[provinceCode].(itemStruct).Cld[cityCode].(itemStruct).Cld[districtCode]
				if ok {
					logrus.Error("repeat item:", item)
					continue
				} else {
					regionMap[provinceCode].(itemStruct).Cld[cityCode].(itemStruct).Cld[districtCode] = itemStruct{
						Code: item.Code,
						Name: item.Name,
					}
				}
			}
		}

		//处理部分县级升地级的情况
		for k, v := range regionMap {
			itemProvince := v.(itemStruct)
			for kk, vv := range itemProvince.Cld {
				itemCity := vv.(itemStruct)
				if itemCity.Name == "" {
					for kkk, vvv := range itemCity.Cld {
						itemDistrict := vvv.(itemStruct)
						itemProvince.Cld[kkk] = itemDistrict
					}
					delete(itemProvince.Cld, kk)
					regionMap[k] = itemProvince
				}
			}
		}

		regionBytes, _ := jsoniter.Marshal(regionMap)
		redisValue := string(regionBytes)
		common.Redis.Set(common.RedisKeyRegion, redisValue, common.RedisTtlRegion)
	} else {
		err := jsoniter.Unmarshal([]byte(redisValue), &regionMap)
		if err != nil {
			utils.ResErrDefault(c)
			return
		}
	}

	result := map[string]interface{}{"region": regionMap}
	utils.ResSuc(c, result)
}
