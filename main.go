package main

import (
	_ "suask/internal/packed"

	_ "suask/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/v2/os/gcfg"

	"suask/internal/cmd"

	_ "suask/internal/enum"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
