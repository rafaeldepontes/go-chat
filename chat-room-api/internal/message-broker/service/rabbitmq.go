package service

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	messagebroker "github.com/rafaeldepontes/go-chat/internal/message-broker"

	"github.com/rafaeldepontes/go-chat/pkg/message-broker/rabbitmq"
)

type msgBrokerSvc struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewMsgService() messagebroker.MsgBroker {
	return &msgBrokerSvc{
		conn:    rabbitmq.GetConn(),
		channel: rabbitmq.GetChannel(),
		queue:   rabbitmq.GetQueue(),
	}
}

func (m *msgBrokerSvc) Send(msg []byte) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := m.channel.PublishWithContext(
		ctx,
		"",
		m.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	); err != nil {
		return err
	}
	return nil
}
