package api

import (
	"encoding/json"
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
	"time"
	"unsafe"
)


type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Common struct {
	CreateAt int    `json:"create_at"`
	Id       string `json:"id"`
}

var commonExample = Common{
	CreateAt: int(time.Now().Unix()),
	Id:       "xxx-xxx-xx",
}

type Page struct {
	PageNum  int `json:"page_num" from:"path" maximum:"1000" desc:"页码"`
	PageSize int `json:"page_size" from:"path" maximum:"10000" minimum:"0" desc:"页大小"`
}

func or (ss ...string)string{
	for _, s := range ss {
		if s != ""{
			return s
		}
	}
	return ""
}


func newErrorResp(code int,msg string)*Resp{
	return &Resp{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}

func newSuccessResp(data interface{})(int ,*Resp){
	return 200,&Resp{Code: 0,Data: data}
}

func newErrorHttpResp(code int,err error)(int ,*Resp){
	return mapCodeToHttp(code),&Resp{Code: code,Message: err.Error()}
}



var (
	notFoundError = newErrorResp(404,"not found data")
	deleteSuccessResp = &Resp{Message: "success"}
)


type face struct {
	t,d unsafe.Pointer
}

func isNil(i interface{})bool{
	p := (*face)(unsafe.Pointer(&i))
	return p.t == nil || p.d == nil
}

func configJsonToString(c interface{}) models.JsonConfigString {
	if isNil(c) {
		return ""
	}
	bs, _ := json.Marshal(c)
	return models.JsonConfigString(bs)
}

func convertErrorf(e error, f string, args ...interface{}) (code int, err error) {
	if e == dao.NoElemError {
		return CodeNotFound, fmt.Errorf(f, args...)
	}
	return CodeDbError, fmt.Errorf(f, args...)
}

