package send_code

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/util/grand"
	"github.com/karim-w/go-azure-communication-services/emails"
)

var endpoint = g.Cfg().MustGet(context.TODO(), "email.endpoint").String()
var accessKey = g.Cfg().MustGet(context.TODO(), "email.accessKey").String()
var senderAddr = g.Cfg().MustGet(context.TODO(), "email.senderAddr").String()

func SendCode(email string) (code string, err error) {
	code = generateCode()

	client := emails.NewClient(endpoint, accessKey, nil)
	payload := emails.Payload{
		SenderAddress: senderAddr,
		Content: emails.Content{
			Subject:   "SuAsk注册 - 验证码",
			PlainText: code,
		},
		Recipients: emails.Recipients{
			To: []emails.ReplyTo{
				{
					Address: email,
				},
			},
		},
	}

	_, err = client.SendEmail(context.TODO(), payload)
	return
}

func generateCode() (code string) {
	return grand.Digits(6)
}
