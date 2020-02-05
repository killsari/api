package utils

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

const (
	errCodeKey    = "errCode"
	errMsgKey     = "errMsg"
	errCodeValSuc = uint64(0)
	errMsgValSuc  = "ok"
	errCodeValErr = uint64(http.StatusInternalServerError)
	errMsgValErr  = "error"
)

func ResSucWithoutData(c *gin.Context) {
	ResSuc(c, nil)
}

func ResSuc(c *gin.Context, data interface{}) {
	dataP := &map[string]interface{}{}
	if data != nil {
		bytes, _ := jsoniter.Marshal(data)
		_ = jsoniter.Unmarshal(bytes, dataP)
	}
	(*dataP)[errCodeKey] = errCodeValSuc
	(*dataP)[errMsgKey] = errMsgValSuc
	res(c, dataP)
}

func ResErrDefault(c *gin.Context) {
	ResErrWithMsg(c, errMsgValErr)
}

func ResErrWithMsg(c *gin.Context, msg string) {
	ResErr(c, errCodeValErr, msg)
}

func ResErr(c *gin.Context, code uint64, msg string) {
	ResErrWithData(c, code, msg, nil)
}

func ResErrWithData(c *gin.Context, code uint64, msg string, data interface{}) {
	dataP := &map[string]interface{}{}
	if data != nil {
		bytes, _ := jsoniter.Marshal(data)
		_ = jsoniter.Unmarshal(bytes, dataP)
	}
	(*dataP)[errCodeKey] = code
	(*dataP)[errMsgKey] = msg
	res(c, dataP)
}

func res(c *gin.Context, data interface{}) {
	c.SecureJSON(http.StatusOK, data)
}

type ResPagination struct {
	Current  uint64 `json:"current"`
	PageSize uint64 `json:"pageSize"`
	Total    uint64 `json:"total"`
}

type ResPage struct {
	List       []interface{} `json:"list"`
	Pagination ResPagination `json:"pagination"`
}

func ResPageFormat(list *[]interface{}, page uint64, size uint64, total uint64) *ResPage {
	return &ResPage{List: *list, Pagination: ResPagination{Current: page, PageSize: size, Total: total}}
}

type ResFieldFlow struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Left  uint64 `json:"left"`
	Total uint64 `json:"total"`
}

type ResFlowPage struct {
	List      []interface{} `json:"list"`
	FieldFlow ResFieldFlow  `json:"fieldFlow"`
}

func ResFlowFormat(list *[]interface{}, first string, last string, left uint64, total uint64) *ResFlowPage {
	return &ResFlowPage{List: *list, FieldFlow: ResFieldFlow{First: first, Last: last, Left: left, Total: total}}
}
