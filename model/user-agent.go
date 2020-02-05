package model

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserAgentModel struct {
}

/*
zh|app|0.0.1|640x960|ios12.4.0|134.123,26.43
*/
type UserAgentEntity struct {
	Language   string //zh/en
	Source     string //pc(电脑页)/app/h5/wxh5(微信h5)/wxapp(小程序)/...
	Version    string //0.0.1
	Resolution string //640x960/750x1334/1080x1920
	Device     string //ios12.3.1/ios12.3.2/android7.0/chrome/ie/...
	Location   string //134.123,26.43
	Uid        uint64
	Token      string
}

const (
	UserAgentHeader = "Killsari-UA"
	UserAgentSep    = "|"
)

const (
	UserAgentSourcePc    = "pc"
	UserAgentSourceApp   = "app"
	UserAgentSourceH5    = "h5"
	UserAgentSourceWxH5  = "wxh5"
	UserAgentSourceWeApp = "weapp"
)

const (
	UserAgentDevicePreIOS     = "ios"
	UserAgentDevicePreAndroid = "android"
)

const (
	UserAgentLanguageZh = "zh"
	UserAgentLanguageEn = "en"
)

func (m UserAgentModel) GetUA(c *gin.Context) *UserAgentEntity {
	ua := UserAgentEntity{}
	value, exists := c.Get("ua")
	if exists {
		uaP, ok := value.(*UserAgentEntity)
		if ok {
			return uaP
		} else {
			logrus.Error("ua format error")
		}
	}
	return &ua
}

func (m UserAgentModel) GetUid(c *gin.Context) uint64 {
	return m.GetUA(c).Uid
}

func (m UserAgentModel) GetToken(c *gin.Context) string {
	return m.GetUA(c).Token
}

func (m UserAgentModel) GetLanguage(c *gin.Context) string {
	return m.GetUA(c).Language
}
