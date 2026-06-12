package middleware

import (
	"net"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// campusSubnets 定义校园网网段列表。
// 请根据中大实际校园网 IP 段进行配置，以下为示例。
// 也可以从配置文件读取。
var campusSubnets []*net.IPNet

func init() {
	// 从配置文件读取校园网网段，如果没有配置则使用默认值
	defaultCIDRs := []string{
		"172.27.0.0/16",   // 中大校园网内网段
		"58.249.112.0/24", // 中大校园网公网出口段
		"127.0.0.0/8",     // 本地开发用
	}
	for _, cidr := range defaultCIDRs {
		_, subnet, err := net.ParseCIDR(cidr)
		if err == nil {
			campusSubnets = append(campusSubnets, subnet)
		}
	}
}

// InitCampusSubnets 从配置文件加载校园网网段（在 server 启动时调用）。
func InitCampusSubnets(ctx interface{}) {
	cidrs := g.Cfg().MustGet(nil, "campus.subnets").Strings()
	if len(cidrs) == 0 {
		return
	}
	campusSubnets = nil
	for _, cidr := range cidrs {
		_, subnet, err := net.ParseCIDR(cidr)
		if err == nil {
			campusSubnets = append(campusSubnets, subnet)
		}
	}
}

// IsCampusNetwork 判断给定 IP 是否在校园网网段内。
func IsCampusNetwork(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, subnet := range campusSubnets {
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}

// GetClientIP 从请求中获取客户端真实 IP。
// 优先从 X-Forwarded-For / X-Real-IP 头获取（适配反向代理场景）。
func GetClientIP(r *ghttp.Request) string {
	// X-Forwarded-For 可能包含多个 IP，取第一个
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	// 回退到 RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// CampusNetworkCheck 中间件：校验请求是否来自校园网。
func CampusNetworkCheck(r *ghttp.Request) {
	clientIP := GetClientIP(r)
	if !IsCampusNetwork(clientIP) {
		g.Log().Infof(r.Context(), "非校园网访问被拒绝, IP: %s", clientIP)
		r.Response.WriteJsonExit(g.Map{
			"code":    403,
			"message": "请在校园网环境下访问",
		})
		return
	}
	r.Middleware.Next()
}
