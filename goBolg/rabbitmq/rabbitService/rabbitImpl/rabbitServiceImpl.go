package rabbitImpl

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"goBolg/dto"
	"goBolg/rabbitmq/rabbitService"
)

type rabbitServiceImpl struct {
	connection *amqp.Connection
}

// NewRabbitService 创建新的 RabbitService 实例
func NewRabbitService(conn *amqp.Connection) rabbitService.RabbitService {
	return &rabbitServiceImpl{connection: conn}
}

func (r *rabbitServiceImpl) SendEmail(emailDTO dto.EmailDTO) error {
	channel, err := r.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// 序列化消息
	body, err := json.Marshal(emailDTO)
	if err != nil {
		return err
	}

	// 发送消息到 RabbitMQ
	err = channel.Publish(
		"EMAIL_EXCHANGE", // 交换机
		"",               // 路由键
		false,            // Mandatory
		false,            // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
