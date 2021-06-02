package main

import (
	"bufio"
	"encoding/json"
	"filestore-server/config"
	dbplayer "filestore-server/db"
	"filestore-server/mq"
	"filestore-server/service/transfer/process"
	_ "github.com/micro/go-plugins/registry/consul/v2"

	"filestore-server/store/oss"
	"fmt"
	"github.com/micro/go-micro/v2"
	"log"
	"os"
	"time"
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
	service := micro.NewService(
		micro.Name("go.micro.service.transfer"), // 服务名称
		micro.RegisterTTL(time.Second*10),       // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5),   // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
	)
	service.Init()
	// 初始化mq client


	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startTranserService() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")
	mq.Init()
	mq.StartConsumer(
		config.TransOSSQueueName,
		"transfer_oss",
		process.Transfer)
}
func main(){
	go startTranserService()
	startRPCService()
}