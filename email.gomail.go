package emailSender

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func (s *OneSender) SendGoEmail() error {
	// 创建邮件对象
	mail := gomail.NewMessage()
	mail.SetHeader("From", config.From)
	mail.SetHeader("To", s.To)
	mail.SetHeader("Subject", s.Subject)

	// 判断邮件内容类型
	if isHTML(s.Body) {
		mail.SetHeader("Content-Type", "text/html; charset=UTF-8")
	} else {
		mail.SetHeader("Content-Type", "text/plain; charset=UTF-8")
	}
	mail.SetBody("text/html", s.Body)

	// 设置SMTP客户端并进行身份验证
	dialer := gomail.NewDialer(config.Smtp, config.PortInt, config.From, config.From_code)

	// 发送邮件
	if err := dialer.DialAndSend(mail); err != nil {
		return fmt.Errorf("Error sending email to %s: %v", s.To, err)
	}

	return nil
}

func (s *Senders) SendGoEmail() error {
	// 创建邮件对象
	mail := gomail.NewMessage()
	mail.SetHeader("From", config.From)

	// 依次发送每封邮件
	for _, e := range s.Sends {
		mail.SetHeader("To", e.To)
		mail.SetHeader("Subject", e.Subject)

		// 判断邮件内容类型
		if isHTML(e.Body) {
			mail.SetHeader("Content-Type", "text/html; charset=UTF-8")
		} else {
			mail.SetHeader("Content-Type", "text/plain; charset=UTF-8")
		}
		mail.SetBody("text/html", e.Body)

		// 设置SMTP客户端并进行身份验证
		dialer := gomail.NewDialer(config.Smtp, config.PortInt, config.From, config.From_code)

		// 发送邮件
		if err := dialer.DialAndSend(mail); err != nil {
			// 记录发送失败的邮件与错误信息
			fmt.Printf("Error sending email to %s: %v\n", e.To, err)
			continue
		}
	}

	return nil
}

func (s *SenderSame) SendGoEmail() error {
	// 创建邮件对象
	mail := gomail.NewMessage()
	mail.SetHeader("From", config.From)
	mail.SetHeader("To", s.To...)
	mail.SetHeader("Subject", s.Subject)

	// 判断邮件内容类型
	if isHTML(s.Body) {
		mail.SetHeader("Content-Type", "text/html; charset=UTF-8")
	} else {
		mail.SetHeader("Content-Type", "text/plain; charset=UTF-8")
	}
	mail.SetBody("text/html", s.Body)

	// 设置SMTP客户端并进行身份验证
	dialer := gomail.NewDialer(config.Smtp, config.PortInt, config.From, config.From_code)

	// 发送邮件
	if err := dialer.DialAndSend(mail); err != nil {
		return fmt.Errorf("Error sending email to %s: %v", s.To, err)
	}

	return nil
}
