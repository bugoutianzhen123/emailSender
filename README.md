# EmailSender

## 简介

自用的邮件发送工具

## 如何使用

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

body ：消息体可以是字符串，也可以是html



发送字符串

```go
body := "this is string"
```



发送html

```go
body := "<html>
  <body>
    <h1>欢迎使用邮件系统</h1>
    <p>这是一封 <strong>HTML</strong> 邮件。</p>
  </body>
</html>"
```

