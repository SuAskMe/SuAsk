package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
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
		fmt.Println("---- SuAsk DB 诊断 ----")
		fmt.Printf("CWD:   %s\n", cwd)

		// 让 GoFrame 告诉我们它实际用的是哪个配置文件
		if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
			if filePath, err := adapter.GetFilePath(); err == nil && filePath != "" {
				fmt.Printf("CONFIG FILE: %s\n", filePath)
				if info, err := os.Stat(filePath); err == nil {
					fmt.Printf("CONFIG SIZE: %d bytes\n", info.Size())
				}
			} else {
				fmt.Println("CONFIG FILE: <未找到>，请检查 manifest/config/config.yaml 是否存在且可读")
			}
			fmt.Printf("SEARCH PATHS: %v\n", adapter.GetPaths())
		}

		// 列出所有顶层配置 key（如果连 server/database/... 都没有，说明文件空或解析失败）
		if data, err := g.Cfg().Data(ctx); err == nil {
			keys := make([]string, 0, len(data))
			for k := range data {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			fmt.Printf("TOP-LEVEL KEYS: %v\n", keys)
		}

		link := g.Cfg().MustGet(ctx, "database.default.link").String()
		dbType := g.Cfg().MustGet(ctx, "database.default.type").String()
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
