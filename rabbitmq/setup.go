package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/mateothegreat/go-rabbitmq/consumer"
	"github.com/mateothegreat/go-rabbitmq/management"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/risersh/controller/conf"
)

// Setup sets up the rabbitmq connection and exchanges.
//
// Returns:
//   - err: an error if the setup fails
func Setup() error {
	exchange := management.Exchange{
		Name:    "primary",
		Type:    "topic",
		Durable: true,
		Queues: []management.Queue{
			{
				Name:    "broker",
				Durable: true,
			},
		},
	}

	m := management.Management{}
	err := m.Connect(conf.Config.RabbitMQ.URI, management.SetupArgs{
		Exchanges: []management.Exchange{exchange},
	})
	if err != nil {
		return err
	}

	return nil
}

// StartConsuming starts consuming messages from the broker queue and
// sends them to the channel which is listened to by the caller.
// The channel is a blocking channel so the caller must listen to it.
//
// Returns:
//   - ch: a blocking channel of messages
//   - err: an error if the consumer fails to start
func StartConsuming[T any]() (chan Message[T], error) {

	c := consumer.Consumer{}
	err := c.Connect(conf.Config.RabbitMQ.URI)
	if err != nil {
		return nil, err
	}

	// Channel to send messages to the caller
	ch := make(chan Message[T])

	// Channel to receive messages from the broker
	channel := make(chan *amqp.Delivery)

	go func() {
		for {
			payload := <-channel
			message := &Message[T]{}
			err := json.Unmarshal(payload.Body, message)
			if err != nil {
				log.Fatalf("Unmarshal failed: %v", err)
			}
			go func() {
				ch <- *message
			}()
		}
	}()

	// Start consuming messages from the broker
	err = consumer.Consume(&c, "controller", channel, "controller")
	if err != nil {
		log.Fatalf("Consuming error: %v", err)
	}

	return ch, nil
}
