package emailSender

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"net/smtp"
	"strings"
)

type EmailSender interface {
	SendEmail() error
}

type EmailConfig struct {
	Smtp      string `yaml:"Smtp"`      //smtp 地址
	Port      string `yaml:"Port"`      //端口
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
	config = &EmailConfig{
		Smtp:      e.Smtp,
		Port:      e.Port,
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

func NewOneSender(s *OneSender) EmailSender { return s }

func NewSenders(s *Senders) EmailSender { return s }

func NewSenderSame(s *SenderSame) EmailSender { return s }

func (s *OneSender) SendEmail() error {
	// 使用SMTP认证
	auth := smtp.PlainAuth("", config.From, config.From_code, config.Smtp)
	// 连接到SMTP服务器
	conn, err := smtp.Dial(config.Smtp + ":" + config.Port)
	if err != nil {
		return fmt.Errorf("Error connecting to SMTP server: %v", err)
	}
	defer conn.Close()

	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 确保不跳过证书验证
		ServerName:         config.Smtp,
	}

	// 启动STARTTLS加密
	if err := conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("Error starting TLS: %v", err)
	}

	// 进行身份验证
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("Error authenticating: %v", err)
	}

	if err := s.send(conn); err != nil {
		return fmt.Errorf("Error sending email to %s: %v", s.To, err)
	}

	return nil
}

func (s *Senders) SendEmail() error {
	// 使用SMTP认证
	auth := smtp.PlainAuth("", config.From, config.From_code, config.Smtp)
	// 连接到SMTP服务器
	conn, err := smtp.Dial(config.Smtp + ":" + config.Port)
	if err != nil {
		return fmt.Errorf("Error connecting to SMTP server: %v", err)
	}
	defer conn.Close()

	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 确保不跳过证书验证
		ServerName:         config.Smtp,
	}

	// 启动STARTTLS加密
	if err := conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("Error starting TLS: %v", err)
	}

	// 进行身份验证
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("Error authenticating: %v", err)
	}

	//设置发件人
	if err := conn.Mail(config.From); err != nil {
		return fmt.Errorf("Error setting sender: %v", err)
	}

	for _, e := range s.Sends {
		if err := e.send(conn); err != nil {
			fmt.Errorf("Error sending email to %s: %v", e.To, err)
			//记录发送失败的邮件与错误信息，再异步处理
			fmt.Println(e, err)
			continue
		}
	}

	return nil
}

func (s *SenderSame) SendEmail() error {
	// 设置邮件头部信息
	fromHeader := fmt.Sprintf("From: %s\r\n", config.From)
	toHeader := fmt.Sprintf("To: %s\r\n", strings.Join(s.To, ","))
	subjectHeader := fmt.Sprintf("Subject: %s\r\n", s.Subject)

	// 判断body是否为HTML
	var contentTypeHeader string
	if isHTML(s.Body) {
		contentTypeHeader = "Content-Type: text/html; charset=UTF-8\r\n"
	} else {
		contentTypeHeader = "Content-Type: text/plain; charset=UTF-8\r\n"
	}

	message := fromHeader + toHeader + subjectHeader + contentTypeHeader + "\r\n" + s.Body

	// 使用SMTP认证
	auth := smtp.PlainAuth("", config.From, config.From_code, config.Smtp)

	// 连接到SMTP服务器
	conn, err := smtp.Dial(config.Smtp + ":" + config.Port)
	if err != nil {
		return fmt.Errorf("Error connecting to SMTP server: %v", err)
	}
	defer conn.Close()

	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 确保不跳过证书验证
		ServerName:         config.Smtp,
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
	if err := conn.Mail(config.From); err != nil {
		return fmt.Errorf("Error setting sender: %v", err)
	}
	for _, recipient := range s.To {
		//if recipient == config.From {
		//	//log.Printf("can't send email to self")
		//	continue
		//}
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

func (s *OneSender) send(conn *smtp.Client) error {
	// 设置邮件头部信息
	fromHeader := fmt.Sprintf("From: %s\r\n", config.From)
	toHeader := fmt.Sprintf("To: %s\r\n", s.To)
	subjectHeader := fmt.Sprintf("Subject: %s\r\n", s.Subject)

	// 判断body是否为HTML
	var contentTypeHeader string
	if isHTML(s.Body) {
		contentTypeHeader = "Content-Type: text/html; charset=UTF-8\r\n"
	} else {
		contentTypeHeader = "Content-Type: text/plain; charset=UTF-8\r\n"
	}

	// 构造邮件正文
	message := fromHeader + toHeader + subjectHeader + contentTypeHeader + "\r\n" + s.Body

	//设置发件人
	conn.Mail(config.From)
	//if err := conn.Mail(config.From); err != nil {
	//	return fmt.Errorf("Error setting sender: %v", err)
	//}

	// 设置收件人
	if err := conn.Rcpt(s.To); err != nil {
		return fmt.Errorf("Error setting recipient: %v", err)
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
