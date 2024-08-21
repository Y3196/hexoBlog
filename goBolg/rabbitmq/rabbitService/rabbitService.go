package rabbitService

import (
	"goBolg/dto"
)

// RabbitService 是发送电子邮件的服务接口
type RabbitService interface {
	SendEmail(emailDTO dto.EmailDTO) error
}
