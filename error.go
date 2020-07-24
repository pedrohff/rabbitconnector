package rabbitmqconnector

import "github.com/streadway/amqp"

type ConnectorError struct {
	Cause             amqp.Error
	IsChannelError    bool
	IsConnectionError bool
}

func (c ConnectorError) Error() string {
	return c.Cause.Error()
}
