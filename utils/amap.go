package utils

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/killsari/api/config"
)

type AmapUtils struct {
}

type amapGeocodes struct {
	FormattedAddress string `json:"formatted_address"` //结构化地址信息
	Location         string `json:"location"`          //坐标点:经度，纬度
}

type AmapGeo struct {
	Status   string          `json:"status"`   //返回结果状态值:返回值为 0 或 1，0 表示请求失败；1 表示请求成功
	Info     string          `json:"info"`     //返回状态说明: 当 status 为 0 时，info 会返回具体错误原因，否则返回“OK”
	Count    string          `json:"count"`    //返回结果的个数
	Geocodes *[]amapGeocodes `json:"geocodes"` //地理编码信息列表
}

var amapKey string

func init() {
	amapKey = config.Conf.AmapKey
}

func (AmapUtils) AmapGetGeoByAddress(s string) (address string, location string, err error) {
	url := fmt.Sprintf("%s?address=%s&key=%s", "https://restapi.amap.com/v3/geocode/geo", s, amapKey)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return
	}

	amapGeo := AmapGeo{}
	err = jsoniter.Unmarshal(body, &amapGeo)
	if err != nil {
		logrus.Error(err)
		return
	}
	if amapGeo.Status == "1" {
		var count int
		count, err = strconv.Atoi(amapGeo.Count)
		if err != nil {
			logrus.Error(err)
			return
		}
		if count > 0 && amapGeo.Geocodes != nil && len(*amapGeo.Geocodes) > 0 {
			address = (*amapGeo.Geocodes)[0].FormattedAddress
			location = (*amapGeo.Geocodes)[0].Location
			return
		}
	}
	errmsg := "amap geo api error"
	if amapGeo.Status == "0" {
		errmsg = amapGeo.Info
	}
	logrus.Error(errmsg)
	err = errors.New(errmsg)
	return
}
