package version

import (
	"os"
	_ "unsafe"

	"golang.org/x/net/http/httpproxy"
)

//go:linkname httpDotResetProxyConfig net/http.resetProxyConfig
func httpDotResetProxyConfig()

const (
	envVarHttpProxy     = "HTTP_PROXY"
	envVarHttpsProxy    = "HTTPS_PROXY"
	envVarNoProxy       = "NO_PROXY"
	envVarRequestMethod = "REQUEST_METHOD"
)

func OverrideProxyInEnvironment(proxy string) func() {
	oldConfig := httpproxy.FromEnvironment()
	oldRequestMethod := os.Getenv("REQUEST_METHOD")

	os.Setenv(envVarHttpProxy, proxy)
	os.Setenv(envVarHttpsProxy, proxy)
	httpDotResetProxyConfig()

	return func() {
		os.Setenv(envVarHttpProxy, oldConfig.HTTPProxy)
		os.Setenv(envVarHttpsProxy, oldConfig.HTTPSProxy)
		os.Setenv(envVarNoProxy, oldConfig.NoProxy)
		os.Setenv(envVarRequestMethod, oldRequestMethod)
		httpDotResetProxyConfig()
	}
}
