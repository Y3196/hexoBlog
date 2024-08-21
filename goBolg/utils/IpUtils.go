package utils

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/mssola/user_agent"
)

// GetIPAddress 从请求中获取 IP 地址
func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	if ip == "127.0.0.1" {
		ip = getLocalIP()
	}
	if idx := strings.Index(ip, ","); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// getLocalIP 获取主机的非回环本地 IP 地址
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetIPSource 查询 IP 地址的地理位置
func GetIPSource(ip string) string {
	apiURL := "http://ip-api.com/json/" + url.QueryEscape(ip)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "Unknown location"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Unknown location"
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "Unknown location"
	}
	if loc, ok := result["city"].(string); ok {
		return loc
	}
	return "Unknown location"
}

// GetUserAgent 解析 User-Agent 字符串并返回用户代理详细信息
func GetUserAgent(r *http.Request) *user_agent.UserAgent {
	ua := r.Header.Get("User-Agent")
	return user_agent.New(ua)
}

// GetIPAddressFromContext 从上下文中获取 IP 地址
func GetIPAddressFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value("ipAddress").(string); ok {
		return ip
	}
	return "Unknown IP"
}

// GetUserAgentFromContext 从上下文中获取用户代理信息
func GetUserAgentFromContext(ctx context.Context) *user_agent.UserAgent {
	if ua, ok := ctx.Value("userAgent").(string); ok {
		return user_agent.New(ua)
	}
	return user_agent.New("")
}
