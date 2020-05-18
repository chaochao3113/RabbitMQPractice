package main

import (
	"fmt"
	"imooc-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	dycOne := RabbitMQ.NewRabbitMQRouting("exDyc", "dyc_one")
	dycTwo := RabbitMQ.NewRabbitMQRouting("exDyc", "dyc_two")
	for i := 0; i < 10; i++ {
		dycOne.PublishRouting("Hello dyc one!" + strconv.Itoa(i))
		dycTwo.PublishRouting("Hello dyc two!" + strconv.Itoa(i))
		time.Sleep(time.Second)
		fmt.Println(i)
	}
}
