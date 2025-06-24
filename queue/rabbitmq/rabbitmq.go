package rabbitmq

import "github.com/streadway/amqp"

type RabbitMQ struct {
	channel  *Channel
	exchange *Exchange
	queue    *Queue
	key      string
}

func New(fns ...func(obj *RabbitMQ)) *RabbitMQ {
	obj := &RabbitMQ{
		channel:  &Channel{},
		queue:    &Queue{},
		exchange: &Exchange{},
	}
	obj.channel.Init()
	obj.queue.Init()
	obj.exchange.Init()
	for _, fn := range fns {
		if fn != nil {
			fn(obj)
		}
	}
	obj.channel.Connect()
	obj.exchange.InitWithChannel(obj.channel)
	obj.queue.InitWithChannel(obj.channel)
	return obj
}

func (self *RabbitMQ) QueueBind() {
	self.channel.QueueBind(self.exchange, self.queue, self.key, false, nil)
}

func (self *RabbitMQ) Publish(key string, msg amqp.Publishing) {
	self.channel.Exec(func(ch *amqp.Channel) {
		ch.Publish(self.exchange.Name, key, false, false, msg)
	})
}

func (self *RabbitMQ) PublishByte(key string, msg []byte) {
	self.Publish(key, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
}

func (self *RabbitMQ) Consume(cus1 string, handle func(delivery amqp.Delivery)) {
	msgs, err := self.channel.channel.Consume(
		self.queue.Name,
		cus1,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_log.Printf("ERROR RabbitMQ consume message :%v", err)
	}
	// 默认单routine 工作读取消息
	go func() {
		for d := range msgs {
			handle(d)
		}
		_log.Printf("RabbitMQ consume Exit 0")
	}()

}

func (self *RabbitMQ) Close() {
	self.channel.Close()
}
