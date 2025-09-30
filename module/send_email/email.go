package send_email

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
)

type Notice struct {
	User    string
	Type    string
	Content string
	URL     string
}

const _AUTH_CODE_TMPL = `<html><body><p>您的验证码为：$code$</p></body></html>`
const _NEW_MESSAGE_TMPL = `亲爱的$user$，您有新的消息！
[$type$] $content$
详情请访问以下链接：$url$
`
const M16KB = (1 << 16)

var authCodeTargets = []string{"$code$"}
var messageTargets = []string{"$user$", "$type$", "$content$", "$url$"}

var host = g.Cfg().MustGet(context.TODO(), "email.host").String()
var port = g.Cfg().MustGet(context.TODO(), "email.port").Int()
var username = g.Cfg().MustGet(context.TODO(), "email.username").String()
var password = g.Cfg().MustGet(context.TODO(), "email.password").String()
var msgTmplPath = g.Cfg().MustGet(context.TODO(), "email.messageTmpl").String()
var authCodeTmplPath = g.Cfg().MustGet(context.TODO(), "email.authCodeTmpl").String()
var _authCodeCache *cache
var _messageCache *cache
var dailer *gomail.Dialer

func loadTmplFileWithDefault(path string, targets []string, deflt string) string {
	file, err := os.Open(path)
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return deflt
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return deflt
	}
	if fileInfo.Size() > M16KB {
		g.Log().Error(context.TODO(), errors.New("template file too large"))
		return deflt
	}
	var sb strings.Builder
	sb.Grow(int(fileInfo.Size()))
	n, err := io.Copy(&sb, file)
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return deflt
	}
	if n != fileInfo.Size() {
		g.Log().Error(context.TODO(), errors.New("emplate file size not match"))
		return deflt
	}
	data := sb.String()
	for _, target := range targets {
		if !strings.Contains(data, target) {
			g.Log().Error(context.TODO(), errors.New("template file format error"))
			return deflt
		}
	}
	return data
}

func init() {
	var err error
	var messageTmpl = _NEW_MESSAGE_TMPL
	var authCodeTmpl = _AUTH_CODE_TMPL
	if len(authCodeTmplPath) > 0 {
		authCodeTmpl = loadTmplFileWithDefault(authCodeTmplPath, authCodeTargets, authCodeTmpl)
	}
	if len(msgTmplPath) > 0 {
		messageTmpl = loadTmplFileWithDefault(msgTmplPath, messageTargets, messageTmpl)
	}
	if _authCodeCache, err = formatTmpl(authCodeTmpl, authCodeTargets); err != nil {
		panic(err)
	}
	if _messageCache, err = formatTmpl(messageTmpl, messageTargets); err != nil {
		panic(err)
	}
	dailer = gomail.NewDialer(host, port, username, password)
}

func SendCode(email string) (code string, err error) {
	code = generateCode()
	err = SendEmail(email, "SuAsk - 验证码", getAuthCode(_authCodeCache, code))
	if err != nil {
		return "", err
	}
	return code, nil
}

func SendNotice(email string, notice *Notice) error {
	subject := "📩 你有一条来自 Suask 提问箱的新消息！"
	content := getMessage(_messageCache, notice)
	return SendEmail(email, subject, content)
}

func SendEmail(email string, subject string, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "SuAsk<"+username+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	if err := dailer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func generateCode() (code string) {
	return grand.Digits(6)
}
