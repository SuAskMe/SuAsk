package send_code

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"gopkg.in/gomail.v2"

	"github.com/gogf/gf/v2/util/grand"
)

// var endpoint = g.Cfg().MustGet(context.TODO(), "email.endpoint").String()
// var accessKey = g.Cfg().MustGet(context.TODO(), "email.accessKey").String()
// var senderAddr = g.Cfg().MustGet(context.TODO(), "email.senderAddr").String()
var host = g.Cfg().MustGet(context.TODO(), "email.host").String()
var port = g.Cfg().MustGet(context.TODO(), "email.port").Int()
var username = g.Cfg().MustGet(context.TODO(), "email.username").String()
var password = g.Cfg().MustGet(context.TODO(), "email.password").String()
var tmplPath = g.Cfg().MustGet(context.TODO(), "email.tmplPath").String()

var msgTmpl = `<html><body><p>您的验证码为：%s</p></body></html>`

func InitMail() error {
	file, err := os.Open(tmplPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > 1024*1024 {
		return fmt.Errorf("template file too large")
	}
	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		return err
	}
	target := "%s"
	msgTmpl = string(data)
	if !strings.Contains(msgTmpl, target) {
		msgTmpl = `<html><body><p>您的验证码为：%s</p></body></html>`
		return fmt.Errorf("template file format error")
	}
	return nil
}

// func SendCode(email string) (code string, err error) {
// 	code = generateCode()

// 	client := emails.NewClient(endpoint, accessKey, nil)
// 	payload := emails.Payload{
// 		SenderAddress: senderAddr,
// 		Content: emails.Content{
// 			Subject:   "SuAsk注册 - 验证码",
// 			PlainText: code,
// 		},
// 		Recipients: emails.Recipients{
// 			To: []emails.ReplyTo{
// 				{
// 					Address: email,
// 				},
// 			},
// 		},
// 	}

// 	_, err = client.SendEmail(context.TODO(), payload)
// 	return
// }

func SendCode(email string) (code string, err error) {
	code = generateCode()
	m := gomail.NewMessage()
	m.SetHeader("From", "SuAsk<"+username+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "SuAsk - 验证码")
	m.SetBody("text/html", fmt.Sprintf(msgTmpl, code))
	d := gomail.NewDialer(host, port, username, password)
	if err := d.DialAndSend(m); err != nil {
		return "", err
	}
	return code, nil
}

func generateCode() (code string) {
	return grand.Digits(6)
}
