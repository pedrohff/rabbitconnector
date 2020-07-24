package rabbitmqconnector

import (
	"github.com/streadway/amqp"
	"time"
)

type Config struct {
	Connection     ConnectionConfig
	Consumers      []ConsumerConfigGroup
	RestartTimeout time.Duration
}

type ConsumerConfigGroup struct {
	Exchange            ExchangeConfig
	QueueDeclaration    QueueDeclarationConfig
	QueueBinding        QueueBindingConfig
	ConsumerDeclaration ConsumerDeclarationConfig
	Consumer            Consumer
}

type ConnectionConfig struct {
	Dial   string
	Config amqp.Config
}

type ExchangeConfig struct {
	Name       string
	Kind       string
	NoWait     bool
	Internal   bool
	AutoDelete bool
	Durable    bool
	UseDelay   bool
	Args       amqp.Table
}

type QueueDeclarationConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type QueueBindingConfig struct {
	QueueName    string
	QueueKey     string
	ExchangeName string
	NoWait       bool
	Args         amqp.Table
}

type ConsumerDeclarationConfig struct {
	QueueName string
	HostName  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}
