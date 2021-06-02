package main
import(
	"filestore-server/service/apigw/route"
)
//启用API作为一个网关或代理，来作为微服务访问的单一入口。
//它应该在您的基础架构的边缘运行。它将HTTP请求转换为RPC并转发给相应的服务。
func main() {
	r:=route.Router()
	r.Run(":8080")
}
