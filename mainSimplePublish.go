package main

import (
	"fmt"
	"imooc-rabbitmq/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("dycSimple")
	rabbitmq.PublishSimple("Hello dyc!")
	fmt.Println("发送成功!")

}
