package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// checkDatabase 做启动期自检：
//  1. 打印当前工作目录、DB 配置 link、最终解析出的绝对路径；
//  2. 如果 sqlite 文件路径是相对路径且文件不存在，直接 panic（避免默默建空库）；
//  3. 访问 users 表统计行数，为 0 则报错——几乎可以肯定是用错了数据库文件。
//
// 这个函数会在 main 命令启动时调用，失败即退出，不让"带病运行的空库"蒙混过关。
func checkDatabase(ctx context.Context) error {
	cwd, _ := os.Getwd()
	g.Log().Infof(ctx, "[db-check] CWD=%s", cwd)

	link := g.Cfg().MustGet(ctx, "database.default.link").String()
	dbType := g.Cfg().MustGet(ctx, "database.default.type").String()
	g.Log().Infof(ctx, "[db-check] database.default.type=%s link=%q", dbType, link)

	// 配置为空：几乎可以确定是配置文件没被加载（路径错、--gf.gcfg.file 指错等）
	if dbType == "" && link == "" {
		return gerror.Newf(
			"[db-check] database 配置为空！config.yaml 很可能根本没被加载。\n"+
				"当前工作目录：%s\n"+
				"可能原因：\n"+
				"  1) config.yaml 被 .gitignore 忽略了，CI 没上传（最常见！）\n"+
				"  2) systemd ExecStart 的 --gf.gcfg.file 指向了不存在的文件\n"+
				"  3) manifest/config/config.yaml 部署路径不对\n"+
				"后续若报 'unable to open database file: out of memory (14)'，\n"+
				"通常是 sqlite 被传了空路径，同样是本问题的衍生表现。\n"+
				"排查：\n"+
				"  ls -la %s/manifest/config/config.yaml\n"+
				"  find %s -maxdepth 3 -name config.yaml\n"+
				"  systemctl cat <服务名> | grep -E 'ExecStart|WorkingDirectory'",
			cwd, cwd, cwd)
	}

	// 对 sqlite 进行额外的路径检查
	if dbType == "sqlite" || strings.HasPrefix(strings.ToLower(link), "sqlite") {
		if path := extractSqlitePath(link); path != "" {
			abs, _ := filepath.Abs(path)
			g.Log().Infof(ctx, "[db-check] sqlite file: %s (abs=%s)", path, abs)
			info, err := os.Stat(abs)
			if err != nil {
				if os.IsNotExist(err) {
					return gerror.Newf(
						"[db-check] sqlite 文件不存在：%s\n"+
							"当前工作目录：%s\n"+
							"可能原因：systemd 的 WorkingDirectory 不对，或手动 scp 时路径放错了。\n"+
							"   排查命令： systemctl show <服务名> -p WorkingDirectory",
						abs, cwd)
				}
				return gerror.Wrapf(err, "[db-check] stat %s", abs)
			}
			g.Log().Infof(ctx, "[db-check] sqlite file OK, size=%d bytes, mode=%s",
				info.Size(), info.Mode())
		}
	}

	// 数据层面的自检：users 表至少有 1 行。
	// 选择 users 是因为它是几乎所有业务的入口，空则一切异常。
	cnt, err := g.DB().Model("users").Ctx(ctx).Count()
	if err != nil {
		return gerror.Wrapf(err, "[db-check] 查询 users 表失败")
	}
	if cnt == 0 {
		return gerror.Newf(
			"[db-check] users 表为空！几乎可以确认程序连上的不是迁移好的数据库。\n"+
				"当前工作目录：%s\n"+
				"请确认 sqlite 文件路径（见上文日志）指向的是已经迁移好的那份 .db",
			cwd)
	}
	g.Log().Infof(ctx, "[db-check] users rows=%d, database healthy", cnt)
	return nil
}

// extractSqlitePath 解析 "sqlite::@file(/path/to/x.db)" 这种格式，取出 file(...) 里的路径。
func extractSqlitePath(link string) string {
	i := strings.Index(link, "file(")
	if i < 0 {
		return ""
	}
	j := strings.Index(link[i+5:], ")")
	if j < 0 {
		return ""
	}
	return link[i+5 : i+5+j]
}
