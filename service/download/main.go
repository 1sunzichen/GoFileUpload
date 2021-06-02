package main

import (
	"filestore-server/service/download/route"
	"fmt"
	"github.com/micro/go-micro/v2"
	dlProto "filestore-server/service/download/proto"
	dlRpc "filestore-server/service/download/rpc"
	cfg "filestore-server/service/download/config"
	_ "github.com/micro/go-plugins/registry/consul/v2"

)

func startRPCService(){
	service:=micro.NewService(
		micro.Name("go.micro.service.download"),

		)
	service.Init()
	dlProto.RegisterDownloadServiceHandler(service.Server(),new(dlRpc.Download) )
	if err:=service.Run();err!=nil{
		fmt.Println(err)
	}
}
func startAPIService(){
	router:=route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main(){
	go startRPCService()
	startAPIService()
}