package register

import (
	"context"

	"github.com/gogf/gf/v2/util/grand"
	"github.com/karim-w/go-azure-communication-services/emails"
)

const endpoint = "suask.japan.communication.azure.com"
const accessKey = "Ed1aRMaseJf37LVH4N5j5hrlJHT6NnvLHLKQIN4Cm545KnK7y4wmJQQJ99ALACULyCpLpMphAAAAAZCSljny"
const senderAddr = "donotreply@9dae9223-b138-4207-942b-e15bf030c1cf.azurecomm.net"

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
