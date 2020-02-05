package common

import "time"

//errcode编码
const (
	ErrCodeOk                  = uint64(0)   //0:正常
	ErrCodeSys                 = uint64(500) //500:错误
	ErrCodeAuth                = uint64(401) //401:授权错误
	ErrCodeNeedPhoneVerify     = uint64(600) //600:需要验证电话
	ErrCodeNeedAccountComplete = uint64(601) //601:需要补全信息注册
)

const (
	WebUrlPrefix       = "https://killsari.moreinfo.cn"
	AliyunOssUrlPrefix = "https://img.killsari.moreinfo.cn"
	//自定义分隔符!
	AliyunOssSuffixS40  = "@!s40"  //(avatar small 暂不用)
	AliyunOssSuffixS100 = "@!s100" //school/career
	AliyunOssSuffixS120 = "@!s120" //avatar
	AliyunOssSuffixS150 = "@!s150" //(category small 暂不用)
	AliyunOssSuffixS210 = "@!s210" //category large
	AliyunOssSuffixW345 = "@!w345" //strategy small
	AliyunOssSuffixW690 = "@!w690" //strategy large

	LikeRecordDuration = time.Minute
)
