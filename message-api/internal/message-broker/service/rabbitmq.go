package service

import (
	amqp "github.com/rabbitmq/amqp091-go"
	messagebroker "github.com/rafaeldepontes/go-chat/internal/message-broker"

	"github.com/rafaeldepontes/go-chat/pkg/message-broker/rabbitmq"
)

type msgBrokerSvc struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	queue    *amqp.Queue
	consumer *<-chan amqp.Delivery
}

func NewMsgService() messagebroker.MsgBroker {
	return &msgBrokerSvc{
		conn:     rabbitmq.GetConn(),
		channel:  rabbitmq.GetChannel(),
		queue:    rabbitmq.GetQueue(),
		consumer: rabbitmq.GetConsumer(),
	}
}

func (m *msgBrokerSvc) Process(userChannel *chan []byte) error {
	for data := range *m.consumer {
		*userChannel <- data.Body
	}
	return nil
}
