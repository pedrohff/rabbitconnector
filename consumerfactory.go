package rabbitmqconnector

import (
	rabbitmq_x_delay "github.com/pedrohff/rabbitmq-x-delay"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Consume(<-chan amqp.Delivery)
}

type ConsumerFactory interface {

	// Requires channel
	AddExchange(config ExchangeConfig) error

	// Requires channel
	DeclareQueue(config QueueDeclarationConfig) error

	// Requires channel
	BindQueue(config QueueBindingConfig) error

	// Requires channel
	Consume(consumer Consumer, config ConsumerDeclarationConfig) error
}

func newConsumerFactory(channel *amqp.Channel) ConsumerFactory {
	return &consumerFactory{
		ichannel: channel,
		channel:  channel,
	}
}

type consumerFactory struct {
	ichannel RabbitMQChannel
	channel  *amqp.Channel
}

func (r *consumerFactory) AddExchange(config ExchangeConfig) error {
	var err error
	if config.UseDelay {
		err = rabbitmq_x_delay.CreateDefaultDelayExchange(
			config.Name,
			config.Kind,
			r.channel,
		)
	} else {
		err = r.ichannel.ExchangeDeclare(
			config.Name,
			config.Kind,
			config.Durable,
			config.AutoDelete,
			config.Internal,
			config.NoWait,
			config.Args,
		)
	}

	return err
}

func (r *consumerFactory) DeclareQueue(config QueueDeclarationConfig) error {
	var err error
	_, err = r.ichannel.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)

	return err
}

func (r *consumerFactory) BindQueue(config QueueBindingConfig) error {
	var err error
	err = r.ichannel.QueueBind(
		config.QueueName,
		config.QueueKey,
		config.ExchangeName,
		config.NoWait,
		config.Args,
	)

	return err
}

func (r *consumerFactory) Consume(consumer Consumer, config ConsumerDeclarationConfig) error {
	deliveries, err := r.ichannel.Consume(
		config.QueueName,
		config.HostName,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		return err
	}
	go consumer.Consume(deliveries)
	return nil
}
