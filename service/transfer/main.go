package main

import (
	"bufio"
	"encoding/json"
	"filestore-server/config"
	dbplayer "filestore-server/db"
	"filestore-server/mq"
	"filestore-server/store/oss"
	"fmt"
	"log"
	"os"
)

func ProcessTransfer(msg []byte)bool{
	//1.解析msg
	pubData:=mq.TransferData{

	}
	fmt.Println("qwqwqwqw")
	//字符串解析成json  并复制给 pubData
	err:=json.Unmarshal(msg,&pubData)
	if err!=nil{
		log.Println(err.Error())
		return false
	}

	//2.临时存储路径 创建文件句柄
	filed,err:=os.Open("../../"+pubData.CurLocation)
	//filed,err:=os.Open(pubData.CurLocation)
	if err!=nil{
		fmt.Println("文件读取失败")
		log.Println(err.Error())
		return false
	}
	//3.写入oss
	oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(filed),
		)

	if err!=nil{
		log.Println(err.Error())
		return false
	}
	//4.修改路径为oss路径 更新数据库
	//fmt.Printf("filehash:%s location%s",pubData.FileHash,pubData.CurLocation)
   suc:=dbplayer.UpdateFileLocation(
   	pubData.FileHash,pubData.DestLocation)
   if !suc{
   	return false
   }
   return true
}


func startRPCService() {
	// 初始化mq client
   mq.Init()

	mq.StartConsumer(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}
func main(){
	//go startTranserService()
	startRPCService()
}