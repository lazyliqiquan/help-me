package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"

	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"github.com/lazyliqiquan/help-me/config"
)

// GetMd5 将给定字符串进行md5加密
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GetUUID 生成唯一码
func GetUUID() string {
	return uuid.New().String()
}

// GetRand 生成验证码
func GetRand() string {
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// SendCode 向指定邮箱发送验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = config.Config.SenderMailbox
	e.To = []string{toUserEmail}
	e.Subject = "You received a verification code from help-me debug website !"
	e.HTML = []byte("<b>" + code + "</b>")
	return e.Send(config.Config.SmtpServerPath+":"+config.Config.SmtpServerPort,
		smtp.PlainAuth("", config.Config.SenderMailbox,
			config.Config.SmtpServerVerification, config.Config.SmtpServerPath))
	// return e.SendWithTLS(config.Config.SmtpServerPath+":"+config.Config.SmtpServerPort,
	// 	smtp.PlainAuth("", config.Config.SenderMailbox,
	// 		config.Config.SmtpServerVerification, config.Config.SmtpServerPath),
	// 	&tls.Config{InsecureSkipVerify: true, ServerName: config.Config.SmtpServerPath})
}

func IsNuiStrs(arr ...string) bool {
	for _, e := range arr {
		if e == "" {
			return true
		}
	}
	return false
}
