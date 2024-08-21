package rabbitService

import (
	"crypto/tls"
	"fmt"
	"goBolg/config"
	"gopkg.in/gomail.v2"
	"log"
)

// EmailService 是发送邮件的服务结构体
type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewEmailService 创建一个新的邮件服务
func NewEmailService(config config.MailConfig) *EmailService {
	log.Printf("Creating EmailService with Host: %s, Port: %d", config.Host, config.Port)
	return &EmailService{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		From:     config.From,
	}
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 启用 TLS

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	log.Println("邮件发送成功！")
	return nil
}
