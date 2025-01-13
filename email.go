package emailSender

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"net/smtp"
	"strings"
)

type EmailSender interface {
	SendEmail(to []string, subject string, body string) error
}

type Email struct {
	Smtp      string `yaml:"Smtp"`
	Port      string `yaml:"Port"`
	From      string `yaml:"From"`
	From_code string `yaml:"From_code"`
}

// 手动配置
func InitEmail(smtp string, port string, from string, from_code string) EmailSender {
	return &Email{
		Smtp:      smtp,
		Port:      port,
		From:      from,
		From_code: from_code,
	}
}

// 有viper
func InitEmailWithViper() EmailSender {
	var e Email
	if err := viper.UnmarshalKey("email", &e); err != nil {
		panic(err)
	}
	return &Email{
		Smtp:      e.Smtp,
		Port:      e.Port,
		From:      e.From,
		From_code: e.From_code,
	}
}

func (e *Email) SendEmail(to []string, subject string, body string) error {
	// 设置邮件头部信息
	fromHeader := fmt.Sprintf("From: %s\r\n", e.From)
	toHeader := fmt.Sprintf("To: %s\r\n", strings.Join(to, ","))
	subjectHeader := fmt.Sprintf("Subject: %s\r\n", subject)
	message := fromHeader + toHeader + subjectHeader + "\r\n" + body

	// 使用SMTP认证
	auth := smtp.PlainAuth("", e.From, e.From_code, e.Smtp)

	// 连接到SMTP服务器
	conn, err := smtp.Dial(e.Smtp + ":" + e.Port)
	if err != nil {
		return fmt.Errorf("Error connecting to SMTP server: %v", err)
	}
	defer conn.Close()

	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 确保不跳过证书验证
		ServerName:         e.Smtp,
	}

	// 启动STARTTLS加密
	if err := conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("Error starting TLS: %v", err)
	}

	// 进行身份验证
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("Error authenticating: %v", err)
	}

	// 设置发件人和收件人
	if err := conn.Mail(e.From); err != nil {
		return fmt.Errorf("Error setting sender: %v", err)
	}
	for _, recipient := range to {
		if err := conn.Rcpt(recipient); err != nil {
			return fmt.Errorf("Error setting recipient: %v", err)
		}
	}

	// 发送邮件正文
	wc, err := conn.Data()
	if err != nil {
		return fmt.Errorf("Error getting data writer: %v", err)
	}
	_, err = wc.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("Error sending email body: %v", err)
	}
	wc.Close()

	return nil
}
