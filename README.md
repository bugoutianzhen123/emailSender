# EmailSender

## 简介

自用的邮件发送工具

## 如何使用

#### 配置

配置相关信息

```go
type EmailConfig struct {
	Smtp      string `yaml:"Smtp"`      //smtp 地址
	Port      string `yaml:"Port"`      //端口
	From      string `yaml:"From"`      //用于发送邮件的邮箱
	From_code string `yaml:"From_code"` //授权码
}
```

手动配置

```go
InitEmailConfig(&emailSender{
    Smtp      :"smtp",
    Port      :"port"
    From      :"from"
    From_code :"from_code"
})
```

当程序使用了viper时

```go
InitEmailConfigWithViper()
```

#### 接口描述

```go
type EmailSender interface {
	SendEmail() error
}
```

实现的对象

```go
type OneSender struct {
	To      string
	Subject string
	Body    string
}

type Senders struct {
	Sends []OneSender
}

type SenderSame struct {
	To      []string
	Subject string
	Body    string
}
```

#### 如何调用

发送单个邮件

```go
email := NewOneSender(&OneSender{
     To: "example@example.com",
	 Subject: "Test Subject",
     Body:    "Test Body",
})
email.SendEmail()
```

发送多个邮件

```go
emails := NewSenders(&Senders{
     Sends: []OneSender{
       {
          To: "example1@example1.com",
	      Subject: "Test Subject1",
          Body:    "Test Body1",
       },
       {
          To: "example2@example2.com",
	      Subject: "Test Subject2",
          Body:    "Test Body2",
       },
     }
})
emails.SendEmail()
```

群发单个邮件

```go
emails := NewSenderSame(&SenderSame{
		To: []string{
			"example1@example1.com",
			"example2@example2.com",
			"example3@example3.com",
		},
		Subject: "Test Send Same Text to Different MailBox",
		Body:    "Hello World!",
	})
```

