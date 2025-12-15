package rabbitmq

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn     *amqp.Connection
	ch       *amqp.Channel
	queue    *amqp.Queue
	consumer *<-chan amqp.Delivery
)

// openConn opens the rabbitmq connection without having to take care of the logic,
// just calling the function should give you an error if any.
func openConn() error {
	connection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		fmt.Printf("Error while trying to connect with rabbitMQ: %s\n", err)
		return err
	}
	conn = connection
	return nil
}

// openChannel opens the unique channel to process the bulk of AMQP messages
// without having to take care of the logic, just calling the function should
// give you an error if any.
func openChannel() error {
	channel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error while trying to create the unique channel: %s\n", err)
		return err
	}
	ch = channel
	return nil
}

func openQueue() error {
	q, err := ch.QueueDeclare(
		"message_queue", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		fmt.Printf("Error while trying to create the unique channel: %s\n", err)
		return err
	}
	queue = &q
	return nil
}

func openConsumer() error {
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		fmt.Printf("Error while trying to create the consumer: %s\n", err)
		return err
	}
	consumer = &msgs
	return nil
}

// GetConn retrives the ampq connection.
func GetConn() *amqp.Connection {
	if conn != nil {
		return conn
	}
	_ = openConn()
	return conn
}

// GetChannel retrives the unique channel to process data.
func GetChannel() *amqp.Channel {
	if ch != nil {
		return ch
	}
	_ = openChannel()
	return ch
}

// GetQueue retrives the queue connection.
func GetQueue() *amqp.Queue {
	if queue != nil {
		return queue
	}
	_ = openQueue()
	return queue
}

func GetConsumer() <-chan amqp.Delivery {
	if consumer != nil {
		return *consumer
	}
	_ = openConsumer()
	return *consumer
}

func CloseConn() error {
	return conn.Close()
}

func CloseChan() error {
	return ch.Close()
}
