package rabbitmq

import "github.com/streadway/amqp"

type Queue struct {
	queue       *amqp.Queue
	Name        string
	Persistence bool
	AutoDelete  bool
	Exclusive   bool
	Block       bool
	Args        amqp.Table
}

func (self *Queue) Init() {
	self.Persistence = true
}

func (self *Queue) InitWithChannel(ch *Channel) {
	ch.Exec(func(obj *amqp.Channel) {
		if self.Name == "" {
			_log.Print("ERROR RabbitMQ queue name is empty \n")
			return
		}
		q, err := obj.QueueDeclare(
			self.Name,
			self.Persistence,
			self.AutoDelete,
			self.Exclusive,
			self.Block,
			self.Args,
		)
		if err != nil {
			_log.Printf("ERROR RabbitMQ queue init : %v \n", err)
			return
		}
		self.queue = &q
	})
}
