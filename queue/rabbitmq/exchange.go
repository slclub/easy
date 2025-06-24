package rabbitmq

import "github.com/streadway/amqp"

type Exchange struct {
	Name         string
	ExchangeType string
	Persistence  bool //是否持久化
	AutoDelete   bool // 是否自动删除
	Interactive  bool // true 标识exchange 不可以被推送client消息
	Block        bool // 是否阻塞
	Args         amqp.Table
}

func (self *Exchange) Init() {
	self.Persistence = true
}

func (self *Exchange) InitWithChannel(ch *Channel) {
	ch.Exec(func(obj *amqp.Channel) {
		err := obj.ExchangeDeclare(
			self.Name,
			self.ExchangeType,
			self.Persistence,
			self.AutoDelete,
			self.Interactive,
			self.Block,
			self.Args,
		)
		if err != nil {
			_log.Printf("ERROR RabbitMQ exchange init : %v\n", err)
		}

	})
}
