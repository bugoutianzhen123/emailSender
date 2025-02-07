package emailSender

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

type EmailSender interface {
	SendEmail() error
}

type EmailConfig struct {
	Smtp      string `yaml:"Smtp"` //smtp 地址
	Port      string `yaml:"Port"` //端口
	PortInt   int    `yaml:"PortInt"`
	From      string `yaml:"From"`      //用于发送邮件的邮箱
	From_code string `yaml:"From_code"` //授权码
}

var config *EmailConfig

func InitEmailConfig(emailConfig *EmailConfig) error {
	err := CheckConfig(emailConfig)
	if err != nil {
		return err
	}
	config = emailConfig
	return nil
}

// 当程序利用viper读取配置文件时，可以利用此函数快速船舰
// email:
//
//	Smtp    : "smtp.email.com"
//	Port    : "port"
//	From      : "email@example.com"
//	From_code : "authorize_code"
func InitEmailConfigWithViper() error {
	var e EmailConfig
	if err := viper.UnmarshalKey("email", &e); err != nil {
		panic(err)
	}
	port, _ := strconv.Atoi(e.Port)
	config = &EmailConfig{
		Smtp:      e.Smtp,
		Port:      e.Port,
		PortInt:   port,
		From:      e.From,
		From_code: e.From_code,
	}
	err := CheckConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func CheckConfig(emailConfig *EmailConfig) error {
	// 确保传入的配置对象不为 nil
	if emailConfig == nil {
		return fmt.Errorf("email configuration cannot be nil")
	}

	// 验证配置项是否有效
	if emailConfig.Smtp == "" {
		return fmt.Errorf("SMTP address is required")
	}
	if emailConfig.Port == "" {
		return fmt.Errorf("Port is required")
	}
	if emailConfig.From == "" {
		return fmt.Errorf("From email address is required")
	}
	if emailConfig.From_code == "" {
		return fmt.Errorf("From email authorization code is required")
	}
	return nil
}

type OneSender struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Senders struct {
	Sends []OneSender `json:"sends"`
}

type SenderSame struct {
	To      []string `json:"tos"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func (s *OneSender) SendEmail() error {
	switch config.Smtp {
	case "smtp.qq.com":
		return s.SendGoEmail()
	default:
		return s.SendSmtpEmail()
	}
}

func (s *Senders) SendEmail() error {
	switch config.Smtp {
	case "smtp.qq.com":
		return s.SendGoEmail()
	default:
		return s.SendSmtpEmail()
	}
}

func (s *SenderSame) SendEmail() error {
	switch config.Smtp {
	case "smtp.qq.com":
		return s.SendGoEmail()
	default:
		return s.SendSmtpEmail()
	}
}

func NewOneSender(s *OneSender) EmailSender { return s }

func NewSenders(s *Senders) EmailSender { return s }

func NewSenderSame(s *SenderSame) EmailSender { return s }
