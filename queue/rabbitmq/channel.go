package rabbitmq

import (
	"github.com/slclub/go-tips/logf"
	"github.com/streadway/amqp"
)

var _log logf.Logger

func init() {
	if _log == nil {
		_log = logf.New()
	}
}

func Log(log logf.Logger) {
	_log = log
}

// rabbit mq channel
type Channel struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue

	// 链接信息
	MQURL string
}

func (self *Channel) Init() {

}
func (self *Channel) Connect() {
	conn, err := amqp.Dial(self.MQURL)
	if err != nil {
		_log.Printf("ERROR RabbitMQ Failed to connect : %v\n", err)
		return
	}
	self.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		_log.Printf("ERROR RabbitMQ Failed to open a channel : %v\n", err)
		return
	}
	self.channel = ch
}

func (self *Channel) Close() {
	if self.channel != nil {
		self.channel.Close()
	}
	if self.conn != nil {
		self.conn.Close()
	}
}

func (self *Channel) Exec(fns ...func(channel *amqp.Channel)) {
	for _, fn := range fns {
		if fn != nil {
			fn(self.channel)
		}
	}
}

func (self *Channel) QueueBind(ex *Exchange, q *Queue, key string, noWait bool, args amqp.Table) {
	exname := ""
	if ex != nil {
		exname = ex.Name
	}
	qname := ""
	if q != nil {
		qname = q.Name
	}
	err := self.channel.QueueBind(
		qname,
		key,
		exname,
		noWait,
		args,
	)
	if err != nil {
		_log.Printf("ERROR RabbitMQ queueBind :%v \n", err)
	}
}
