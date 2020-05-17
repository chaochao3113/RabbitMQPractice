package main

import (
	"fmt"
	"imooc-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

//工作模式  一个消息只能被一个消费者消费  适用于生产者生产速度大于消费者消费速度  通过多开消费端实现

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("dycSimple")

	for i := 0; i <= 100 ; i++ {
		rabbitmq.PublishSimple("Hello dyc!" + strconv.Itoa(i))
		time.Sleep(time.Second)
		fmt.Println(i)
	}
}
