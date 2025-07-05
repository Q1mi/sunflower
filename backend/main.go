package main

import (
	_ "backend/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2" // 导入MySQL驱动
	"github.com/gogf/gf/v2/os/gctx"

	"backend/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
