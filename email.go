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

type email struct {
	Smtp      string `yaml:"Smtp"`      //smtp 地址
	Port      string `yaml:"Port"`      //端口
	From      string `yaml:"From"`      //用于发送邮件的邮箱
	From_code string `yaml:"From_code"` //授权码
}

// 手动配置发送邮件所需的信息
func InitEmail(smtp string, port string, from string, from_code string) EmailSender {
	return &email{
		Smtp:      smtp,
		Port:      port,
		From:      from,
		From_code: from_code,
	}
}

// 当程序利用viper读取配置文件时，可以利用此函数快速船舰
// email:
//
//	Smtp    : "smtp.email.com"
//	Port    : "port"
//	From      : "email@example.com"
//	From_code : "authorize_code"
func InitEmailWithViper() EmailSender {
	var e email
	if err := viper.UnmarshalKey("email", &e); err != nil {
		panic(err)
	}
	return &email{
		Smtp:      e.Smtp,
		Port:      e.Port,
		From:      e.From,
		From_code: e.From_code,
	}
}

// 发送邮件
// to: 目标邮箱，
// subject: 主题
// body: 内容，可以是string，也可以是完整等等html
func (e *email) SendEmail(to []string, subject string, body string) error {
	// 设置邮件头部信息
	fromHeader := fmt.Sprintf("From: %s\r\n", e.From)
	toHeader := fmt.Sprintf("To: %s\r\n", strings.Join(to, ","))
	subjectHeader := fmt.Sprintf("Subject: %s\r\n", subject)

	// 判断body是否为HTML
	var contentTypeHeader string
	if isHTML(body) {
		contentTypeHeader = "Content-Type: text/html; charset=UTF-8\r\n"
		fmt.Println("Content-Type: text/html; charset=UTF-8")
	} else {
		contentTypeHeader = "Content-Type: text/plain; charset=UTF-8\r\n"
		fmt.Println("Content-Type: text/plain; charset=UTF-8")
	}

	message := fromHeader + toHeader + subjectHeader + contentTypeHeader + "\r\n" + body

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

func isHTML(body string) bool {
	// 简单检查是否包含 HTML 或 SVG 标签
	return strings.Contains(body, "<html>") || strings.Contains(body, "<body>") ||
		strings.Contains(body, "<p>") || strings.Contains(body, "<h1>") ||
		strings.Contains(body, "<svg>") || strings.Contains(body, "<path>")
}
