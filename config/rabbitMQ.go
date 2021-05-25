package config

const (
	AsyncTransferEnable=true //同步异步
	RabbitURL="amqp://root:root@1.116.165.233:5672/"
	TransExchangeName="uploadserver"  //发布名
	TransOSSQueueName="uploadsercer.oss" //队列名
	TransOSSErrQueueName="UploadServer.oss.err"
	TransOSSRoutingKey="oss" //路由匹配的key

)