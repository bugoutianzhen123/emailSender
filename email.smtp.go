package emailSender

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

func (s *OneSender) SendSmtpEmail() error {
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

func (s *Senders) SendSmtpEmail() error {
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

	for _, e := range s.Sends {
		if err := e.send(conn); err != nil {
			fmt.Errorf("Error sending email to %s: %v", e.To, err)
			//记录发送失败的邮件与错误信息，再异步处理
			fmt.Println(e, err)
			continue
		}
		if err := conn.Reset(); err != nil {
			return fmt.Errorf("Error resetting connection: %v", err)
		}
	}

	return nil
}

func (s *SenderSame) SendSmtpEmail() error {
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
	if err := conn.Mail(config.From); err != nil {
		if strings.Contains(err.Error(), "250") {
			fmt.Println("Warning: RCPT TO returned 250 OK (ignored)")
		} else {
			return fmt.Errorf("Error setting sender: %v", err)
		}
	}

	// 设置收件人
	if err := conn.Rcpt(s.To); err != nil {
		if strings.Contains(err.Error(), "250") {
			fmt.Println("Warning: RCPT TO returned 250 OK (ignored)")
		} else {
			return fmt.Errorf("Error setting sender: %v", err)
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
