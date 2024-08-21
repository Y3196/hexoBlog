package rabbitmq

import (
	"github.com/streadway/amqp"
	"goBolg/config"
	"log"
)

type RabbitMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	config     *config.RabbitMQConfig
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQClient{
		connection: conn,
		channel:    ch,
	}, nil
}

func (c *RabbitMQClient) Publish(exchange, key string, body []byte) error {
	return c.channel.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *RabbitMQClient) Close() {
	if err := c.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := c.connection.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}
