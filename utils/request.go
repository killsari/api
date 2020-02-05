package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ReqQueryPage(c *gin.Context) (currentPage uint64, pageSize uint64, err error) {
	pageParam := c.DefaultQuery("currentPage", "1")
	sizeParam := c.DefaultQuery("pageSize", "10")
	currentPage, err = strconv.ParseUint(pageParam, 10, 64)
	if err != nil {
		logrus.Error(err)
		ResErrWithMsg(c, "currentPage must be uint")
		return
	} else if currentPage == 0 {
		ResErrWithMsg(c, "currentPage must not be 0")
		return
	}
	pageSize, err = strconv.ParseUint(sizeParam, 10, 64)
	if err != nil {
		logrus.Error(err)
		ResErrWithMsg(c, "pageSize must be uint")
		return
	} else if pageSize == 0 {
		ResErrWithMsg(c, "pageSize must not be 0")
		return
	}
	return
}

/**
绑定bodyJson到struct
*/
func ReqBindBody(c *gin.Context, reqInfo interface{}) (err error) {
	if err = c.ShouldBind(reqInfo); err != nil {
		logrus.Error(err)
		ResErrWithMsg(c, "input format error")
	}
	return
}

func ReqQueryFlow(c *gin.Context) (first string, last string, count uint64, err error) {
	first = c.DefaultQuery("first", "")
	last = c.DefaultQuery("last", "")
	countParam := c.DefaultQuery("count", "10")
	count, err = strconv.ParseUint(countParam, 10, 64)
	if err != nil {
		logrus.Error(err)
		ResErrWithMsg(c, "count must be uint")
		return
	} else if count == 0 {
		ResErrWithMsg(c, "count must not be 0")
		return
	}
	return
}
