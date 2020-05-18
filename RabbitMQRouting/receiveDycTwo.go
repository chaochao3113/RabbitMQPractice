package main

import "imooc-rabbitmq/RabbitMQ"

func main() {
	dycTwo := RabbitMQ.NewRabbitMQRouting("exDyc", "dyc_two")
	dycTwo.ReceiveRouting()
}
