package  main

import (
	"filestore-server/config"
	"filestore-server/service/upload/route"
	"github.com/micro/go-micro/v2"
	upProto "filestore-server/service/upload/proto"
	upRpc "filestore-server/service/upload/rpc"
	_ "github.com/micro/go-plugins/registry/consul/v2"

	"log"
)
func startRpcService(){
	service:=micro.NewService(micro.Name("go.micro.service.upload"))
	service.Init()
	upProto.RegisterUploadServiceHandler(service.Server(),new(upRpc.Upload))
	if err:=service.Run();err!=nil{
		log.Println(err)
	}
}
func startApiService(){
	router:=route.Route()
	router.Run(config.UploadServiceHost)
}
func main(){
   go startRpcService()
   startApiService()
}
