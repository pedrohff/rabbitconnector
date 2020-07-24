package rabbitmqconnector

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type RabbitMQConnector interface {
	StartListening()
}

func New(config Config, errorNotifier chan error) RabbitMQConnector {
	return &rabbitMQConnector{
		config:         config,
		errorNotifier:  errorNotifier,
		isReconnecting: false,
		reconnectMutex: sync.Mutex{},
	}
}

type rabbitMQConnector struct {
	config         Config
	connection     *amqp.Connection
	iconnection    RabbitMQConnection
	errorNotifier  chan error
	isReconnecting bool
	reconnectMutex sync.Mutex
}

func (r *rabbitMQConnector) createConnection() error {
	var err error
	r.connection, err = amqp.DialConfig(
		r.config.Connection.Dial,
		r.config.Connection.Config,
	)
	r.iconnection = r.connection

	if err != nil {
		r.handleError("could not open connection", err)
		return err
	}

	// Listens for closed connections
	connectionCloseListener := make(chan *amqp.Error)
	r.connection.NotifyClose(connectionCloseListener)
	r.closeListener("connection closed", connectionCloseListener)
	return err
}

func (r *rabbitMQConnector) closeListener(errMsg string, listener chan *amqp.Error) chan *amqp.Error {
	go func() {
		amqpError, ok := <-listener
		if ok && amqpError != nil {
			r.handleError(errMsg, amqpError)
		}
	}()
	return listener
}

func (r *rabbitMQConnector) createChannel() (*amqp.Channel, error) {
	channel, err := r.iconnection.Channel()
	if err != nil {
		r.handleError("could not open channel", err)
		return nil, err
	}

	// Listens for closed channels
	channelCloseListener := make(chan *amqp.Error)
	channel.NotifyClose(channelCloseListener)
	r.closeListener("channel closed", channelCloseListener)
	return channel, nil
}

func (r *rabbitMQConnector) handleError(msg string, err error) {
	if r.errorNotifier != nil {
		r.errorNotifier <- errors.New(fmt.Sprintf("%s : %v", msg, err))
	}

	if !r.isReconnecting {
		r.setReconnecting(true)
		defer r.setReconnecting(false)
		r.closeConnection()

		time.Sleep(r.config.RestartTimeout)

		r.StartListening()
	}
}

func (r *rabbitMQConnector) setReconnecting(b bool) {
	r.reconnectMutex.Lock()
	defer r.reconnectMutex.Unlock()
	r.isReconnecting = b
}

func (r *rabbitMQConnector) closeConnection() {
	if r.iconnection != nil {
		_ = r.iconnection.Close()
	}
}

func (r *rabbitMQConnector) StartListening() {
	err := r.createConnection()
	if err != nil {
		return
	}

	for _, consumerConfig := range r.config.Consumers {
		channel, err := r.createChannel()
		if err != nil {
			return
		}
		factory := newConsumerFactory(channel)

		err = factory.AddExchange(consumerConfig.Exchange)
		if err != nil {
			r.handleError("could not declare exchange", err)
		}

		err = factory.DeclareQueue(consumerConfig.QueueDeclaration)
		if err != nil {
			r.handleError("could not declare queue", err)
		}

		err = factory.BindQueue(consumerConfig.QueueBinding)
		if err != nil {
			r.handleError("could not bind queue to exchange", err)
		}

		err = factory.Consume(consumerConfig.Consumer, consumerConfig.ConsumerDeclaration)
		if err != nil {
			r.handleError("could not start consumer", err)
		}
	}
}
