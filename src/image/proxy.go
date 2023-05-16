package image

import (
	"os"
	_ "unsafe"

	"golang.org/x/net/http/httpproxy"
)

type ProxyConfig struct {
	HTTPProxy     string
	HTTPSProxy    string
	NoProxy       string
	RequestMethod string
}

func GetProxyFromEnvironment() ProxyConfig {
	config := httpproxy.FromEnvironment()
	requestMethod := os.Getenv("REQUEST_METHOD")

	return ProxyConfig{
		HTTPProxy:     config.HTTPProxy,
		HTTPSProxy:    config.HTTPSProxy,
		NoProxy:       config.NoProxy,
		RequestMethod: requestMethod,
	}
}

func SimpleProxyConfig(proxyUrl string) *ProxyConfig {
	return &ProxyConfig{
		HTTPProxy:     proxyUrl,
		HTTPSProxy:    proxyUrl,
		NoProxy:       "",
		RequestMethod: "",
	}
}

func ReconfigureProxyFromEnvironment(proxyConfig *ProxyConfig) func() {
	oldConfig := GetProxyFromEnvironment()

	os.Setenv("HTTP_PROXY", proxyConfig.HTTPProxy)
	os.Setenv("HTTPS_PROXY", proxyConfig.HTTPSProxy)
	os.Setenv("NO_PROXY", proxyConfig.NoProxy)
	os.Setenv("REQUEST_METHOD", proxyConfig.RequestMethod)

	resetCachedProxies()

	return func() {
		os.Setenv("HTTP_PROXY", oldConfig.HTTPProxy)
		os.Setenv("HTTPS_PROXY", oldConfig.HTTPSProxy)
		os.Setenv("NO_PROXY", oldConfig.NoProxy)
		os.Setenv("REQUEST_METHOD", oldConfig.RequestMethod)
	}
}

func applyConfig(config ProxyConfig) {
	os.Setenv("HTTP_PROXY", config.HTTPProxy)
	os.Setenv("HTTPS_PROXY", config.HTTPSProxy)
	os.Setenv("NO_PROXY", config.NoProxy)
	os.Setenv("REQUEST_METHOD", config.RequestMethod)
}

//go:linkname resetCachedProxies net/http.resetProxyConfig
func resetCachedProxies()
