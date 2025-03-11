package main

import (
	_ "suask/internal/packed"
	"suask/utility/send_code"

	_ "suask/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/v2/os/gcfg"

	"suask/internal/cmd"

	_ "suask/internal/enum"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	err := send_code.InitMail()
	if err != nil {
		panic(err)
	}
	cmd.Main.Run(gctx.GetInitCtx())
}
