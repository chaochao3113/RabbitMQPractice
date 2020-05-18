package main

import "imooc-rabbitmq/RabbitMQ"

func main() {
	dycOne := RabbitMQ.NewRabbitMQRouting("exDyc", "dyc_one")
	dycOne.ReceiveRouting()
}
