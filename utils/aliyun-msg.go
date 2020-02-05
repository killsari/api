package utils

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/killsari/api/config"
)

type AliyunMsgUtils struct {
}

var (
	AliyunMsgAccessKeyID     = config.Conf.AliYun.AccessKeyID
	AliyunMsgAccessKeySecret = config.Conf.AliYun.AccessKeySecret
)

const (
	AliyunMsgTemplateCode = "SMS_xxx"
	AliyunMsgSignName     = "xxx"
)

func (u AliyunMsgUtils) SendSms(phone string, code string) (err error) {
	paramByte, _ := jsoniter.Marshal(map[string]string{"code": code})
	param := string(paramByte)

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", AliyunMsgAccessKeyID, AliyunMsgAccessKeySecret)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = AliyunMsgSignName
	request.TemplateCode = AliyunMsgTemplateCode
	request.TemplateParam = param
	_, err = client.SendSms(request)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}
