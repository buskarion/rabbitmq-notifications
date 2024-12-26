package rabbitmq

import "github.com/streadway/amqp"

type MessagingInterface interface {
	Consume(queueName, consumerTag string, autoAck, exclusive, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Publish(exchange, routingKey string, mandatory, immediate bool, msg amqp.Publishing) error
	Close()
}
