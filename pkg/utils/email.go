package utils

import (
	"InnerG/config"
	"InnerG/pkg/errno"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"strconv"
)

// MailSendCode 发送验证码邮件到指定地址。优先使用 config.Smtp 配置，若未初始化则回退到环境变量。
func MailSendCode(to string, code string) error {
	if to == "" {
		return errno.NewErr(errno.InternalServiceErrorCode, "收件人邮箱为空")
	}

	host := config.Smtp.Host
	port := strconv.Itoa(config.Smtp.Port)
	user := config.Smtp.User
	pass := config.Smtp.Password
	from := config.Smtp.From
	fromName := config.Smtp.FromName

	if host == "" || port == "" || user == "" || pass == "" || from == "" {
		return errno.NewErr(errno.InternalServiceErrorCode, "SMTP 配置不完整")
	}

	addr := host + ":" + port

	e := email.NewEmail()
	if fromName != "" {
		e.From = fmt.Sprintf("%s <%s>", fromName, from)
	} else {
		e.From = from
	}
	e.To = []string{to}
	e.Subject = "验证码"
	e.HTML = []byte(fmt.Sprintf("你的验证码为：<h1>%s</h1><p>有效期请以系统设置为准。</p>", code))

	auth := smtp.PlainAuth("", user, pass, host)

	tlsCfg := &tls.Config{ServerName: host}
	if err := e.SendWithTLS(addr, auth, tlsCfg); err == nil {
		return nil
	}

	if err := e.Send(addr, auth); err != nil {
		log.Printf("MailSendCode addr:%s : %w", addr, err.Error())
		return errno.NewErr(errno.InternalServiceErrorCode, "发送邮箱验证码失败")
	}
	return nil
}
