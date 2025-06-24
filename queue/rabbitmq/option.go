package rabbitmq

import "github.com/streadway/amqp"

const (
	EXCHANGE_DIRECT  = "direct"  // 精确匹配路由建 明确指定 队列的场景
	EXCHANGE_FANOUT  = "fanout"  // 发布订阅模式 广播 每个consume 有单独的队列 忽略 路由建
	EXCHANGE_TOPIC   = "topic"   // 灵活 主题模式
	EXCHANGE_HEADERS = "headers" // 复杂路由结构，依据header 健值对 匹配
)

// rabbitmq connection address
func OptionConnect(mqurl string) func(obj *RabbitMQ) {
	return func(this *RabbitMQ) {
		this.channel.MQURL = mqurl
	}
}

func OptionExchangeName(exchangeName string) func(obj *RabbitMQ) {
	return func(this *RabbitMQ) {
		this.exchange.Name = exchangeName
		if this.exchange.ExchangeType == "" {
			this.exchange.ExchangeType = EXCHANGE_DIRECT
		}
	}
}

func OptionQueueName(queueName string) func(obj *RabbitMQ) {
	return func(this *RabbitMQ) {
		this.queue.Name = queueName
	}
}

func OptionPersistence(persist bool) func(obj *RabbitMQ) {
	return func(this *RabbitMQ) {
		this.exchange.Persistence = persist
		this.queue.Persistence = persist
	}
}

func OptionQueueExtendArgs(args amqp.Table) func(obj *RabbitMQ) {
	return func(this *RabbitMQ) {
		this.queue.Args = args
	}
}

func OptionExchangeExtendArgs(args amqp.Table) func(obj *RabbitMQ) {
	return func(this *RabbitMQ) {
		this.exchange.Args = args
	}
}

func OptionRouteKey(key string) func(obj *RabbitMQ) {
	return func(this *RabbitMQ) {
		this.key = key
	}
}
