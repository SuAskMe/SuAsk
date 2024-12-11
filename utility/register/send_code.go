package register

import "github.com/gogf/gf/v2/util/grand"

func SendCode(email string) (code string, err error) {
	code = generateCode()
	// Add email Api Here
	return code, nil
}

func generateCode() (code string) {
	return grand.Digits(6)
}
