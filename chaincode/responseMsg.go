package main

import (
	"encoding/json"
	"fmt"
)

type ResponseMsg struct {
	Code int //相应消息代码（0：成功 1：失败）
	dec string // 消息内容
}

func GetMsgByte(code int, dec string) ([]byte,error) {
	msgByte, err := getMsg(code,dec)
	if err != nil {
		return nil,err
	}
	return msgByte[:],nil
}

func GetMsgString(code int, dec string) (string,error) {
	msgByte, err := getMsg(code,dec)
	if err != nil {
		return "",err
	}
	return string(msgByte),nil
}

func getMsg(code int, dec string) ([]byte, error) {

	var crm ResponseMsg
	crm.Code = code
	crm.dec = dec

	msg,err := json.Marshal(crm)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return msg,nil
}
