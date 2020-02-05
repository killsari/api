package common

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/killsari/api/config"
	"time"
)

var Redis *redis.Client

const (
	RedisKeyPre = "api.killsari:"

	RedisKeyCode = RedisKeyPre + "code_%s_%s" //验证码[code_reg_phone]

	RedisKeyToken = RedisKeyPre + "token_%s" //授权token[token_xxxxxx]

	RedisKeyRegion = RedisKeyPre + "region" //行政区划（省市/地市/县区）

)

const (
	RedisTtlCode = time.Minute //phone验证码有效期

	RedisTtlToken       = time.Hour * 24 * 30 //token有效期
	RedisTtlExtendToken = time.Hour * 24 * 10 //续约时间

	RedisTtlRegion = time.Hour * 24 * 30 //行政区划有效期

)

func init() {
	redisConfig := config.Conf.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})
	Redis = client
}
