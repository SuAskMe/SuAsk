package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

// DBInfo 是一个诊断命令：./main_test dbinfo
// 不启动 HTTP 服务，只输出数据库相关的排错信息，非常适合在部署环境里排查
// "明明数据已迁移，但所有接口都像读了空库" 这类问题。
var DBInfo = gcmd.Command{
	Name:  "dbinfo",
	Usage: "dbinfo",
	Brief: "print database diagnostic info then exit",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		cwd, _ := os.Getwd()
		link := g.Cfg().MustGet(ctx, "database.default.link").String()
		dbType := g.Cfg().MustGet(ctx, "database.default.type").String()

		fmt.Println("---- SuAsk DB 诊断 ----")
		fmt.Printf("CWD:   %s\n", cwd)
		fmt.Printf("TYPE:  %s\n", dbType)
		fmt.Printf("LINK:  %q\n", link)

		if path := extractSqlitePath(link); path != "" {
			abs, _ := filepath.Abs(path)
			fmt.Printf("SQLite file (from link): %s\n", path)
			fmt.Printf("SQLite file abs path:    %s\n", abs)
			if info, err := os.Stat(abs); err == nil {
				fmt.Printf("SQLite file size:        %d bytes\n", info.Size())
			} else {
				fmt.Printf("SQLite file stat error:  %v\n", err)
			}
		}

		// 几张关键业务表的行数，迁移是否成功一眼可见
		tables := []string{"users", "files", "questions", "answers",
			"teachers", "settings", "notifications", "favorites"}
		fmt.Println("\n表行数：")
		for _, t := range tables {
			cnt, err := g.DB().Model(t).Ctx(ctx).Count()
			if err != nil {
				fmt.Printf("  %-15s ERROR: %v\n", t, err)
				continue
			}
			fmt.Printf("  %-15s %d\n", t, cnt)
		}
		fmt.Println("\nOK.")
		return nil
	},
}
