package util

import (
	"encoding/json"
	"log"
)

type RespMsg struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{}`json:"data"`
}

func NewRespMsg(code int,msg string,data interface{})*RespMsg  {
	return &RespMsg{
		Code:code,
		Msg:msg,
		Data:data,
	}
}
//JSONBytes   对象转json格式的string
func (resp *RespMsg)JSONBytes()[]byte{
	r,err:=json.Marshal(resp)
	if err!=nil{
		log.Println(err)
	}
	return r
}

//JSONString 对象转json格式的string
func (resp *RespMsg)JSONString()string{
	r,err:=json.Marshal(resp)
	if err!=nil{
		log.Println(err)
	}
	return string(r)
}

