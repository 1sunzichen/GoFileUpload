package mq
import(
	"filestore-server/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error
// Init : 初始化MQ连接信息
func Init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}
//初始化通道
func initChannel(rabbitHost string)bool{

	//创建通道
	if channel !=nil{
		return true
	}
	//获得连接
	conn,err:=amqp.Dial(rabbitHost)
	if err!=nil{
		log.Println(err.Error())
		fmt.Println("连接失败")
		return false
	}
	//打开通道
	channel,err=conn.Channel()
	if err!=nil{
		log.Println(err.Error())
		return false
	}
	return true
}
//Publish 发布消息
func Publish(exchange,routingKey string,msg []byte)bool{
	//判断channel是否正常

	if !initChannel(config.RabbitURL){
		fmt.Println("mq连接初始化失败")
		return false
	}
	fmt.Printf("mq连接初始化成功%v\n %v\n ",exchange,routingKey)
	//
	err:=channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg,

		})
	if err!=nil{
		log.Println(err.Error())
		return false
	}
	return true

}