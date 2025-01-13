# EmailSender

## 简介

自用的邮件发送工具



手动配置

```go
InitEmail(smtp string, port string, from string, from_code string) EmailSender
```



当程序使用了viper时

```go
InitEmailWithViper() EmailSender
```





接口描述

```go
type EmailSender interface {
	SendEmail(to []string, subject string, body string) error
}
```

to ：目标邮箱，可以是多个

subject： 主题

body ：消息体