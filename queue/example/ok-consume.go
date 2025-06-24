package main

import (
	"fmt"
	"github.com/slclub/easy/queue/rabbitmq"
	"github.com/streadway/amqp"
	"time"
)

var rabbit *rabbitmq.RabbitMQ

func main() {
	initialize()
	rabbit.QueueBind()
	rabbit.Consume("q1", func(res amqp.Delivery) {
		fmt.Printf("q1.read : %v\n", string(res.Body))
	})
	time.Sleep(time.Second * 50)
}

func initialize() {
	rabbit = rabbitmq.New(
		rabbitmq.OptionConnect("amqp://xyj:slclub@192.168.11.44:5672/xyj"),
		rabbitmq.OptionExchangeName("E1"),
		rabbitmq.OptionQueueName("Q2"), // 修改队列名称
	)
}
