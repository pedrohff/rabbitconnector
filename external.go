package rabbitmqconnector

import "github.com/streadway/amqp"

type RabbitMQConnection interface {
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
	Channel() (*amqp.Channel, error)
	Close() error
}

type RabbitMQChannel interface {
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Close() error
}
