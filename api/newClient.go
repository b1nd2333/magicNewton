package api

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"time"
)

func newHTTPClientWithProxy(proxyAddress string) (*http.Client, error) {
	// 解析代理地址
	proxyURL, err := url.Parse(proxyAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy address: %v", err)
	}
	transport := &http.Transport{}

	if proxyURL.Scheme == "socks5" {
		// 设置 SOCKS5 代理并进行身份验证
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("failed to create SOCKS5 dialer: %v", err)
		}
		// 创建 HTTP Transport 使用 SOCKS5 代理
		transport = &http.Transport{
			Dial: dialer.Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略 HTTPS 错误
			},
		}
	} else {
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略 HTTPS 错误
			},
		}
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return client, nil
}
