package main

import (
	_ "suask/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"suask/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
