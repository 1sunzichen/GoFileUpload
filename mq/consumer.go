package mq

import (
	"fmt"
	"log"
)
var done chan bool
//消费者 开始监听队列 获取消息
func StartConsumer(qName,cName string,callback func(msg []byte)bool) {

	//1.通过channel.Consume 获取消息信道
	msgs, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 2 循环获取队列的消息
	done=make(chan bool)
	//

	go func() {
		for msg:=range msgs{
            fmt.Printf("消息:%v  %c",msg,msg.Body)
			suc:=callback(msg.Body)
			if !suc{
				//TODO 写到另一个队列
			}
		}
	}()
    <-done
    //关闭channel
    channel.Close()
}