package main

import (
	"filestore-server/service/account/handler"
	"filestore-server/service/account/proto/proto"
	"github.com/micro/go-micro/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	"log"
	//"time"
)

func main() {
	service:=micro.NewService(
		//go.micro.service.user
		micro.Name("go.micro.service.user"),
		//micro.RegisterTTL(time.Second*10),
		//micro.RegisterInterval(time.Second*5),
		micro.Metadata(map[string]string{"protocol" : "http"}),
		)

	service.Init()

	proto.RegisterUserServiceHandler(service.Server(),new(handler.User))
	if err:=service.Run();err!=nil{
		log.Println(err)
	}
}
