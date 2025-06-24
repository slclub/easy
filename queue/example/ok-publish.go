package main

import (
	"github.com/slclub/easy/queue/rabbitmq"
)

var rabbit *rabbitmq.RabbitMQ

func main() {
	initialize()
	rabbit.QueueBind()
	rabbit.PublishByte("", []byte("ok-test1"))
}

func initialize() {
	rabbit = rabbitmq.New(
		rabbitmq.OptionConnect("amqp://xyj:slclub@192.168.11.44:5672/xyj"),
		rabbitmq.OptionExchangeName("E1"),
		rabbitmq.OptionQueueName("Q1"),
	)
}
