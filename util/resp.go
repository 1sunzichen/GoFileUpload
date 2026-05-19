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
// JSONBytes converts an object to a JSON-format byte slice
func (resp *RespMsg)JSONBytes()[]byte{
	r,err:=json.Marshal(resp)
	if err!=nil{
		log.Println(err)
	}
	return r
}

// JSONString converts an object to a JSON-format string
func (resp *RespMsg)JSONString()string{
	r,err:=json.Marshal(resp)
	if err!=nil{
		log.Println(err)
	}
	return string(r)
}

