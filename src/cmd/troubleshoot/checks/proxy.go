package checks

import (
	"github.com/Dynatrace/dynatrace-operator/src/cmd/troubleshoot"
	"github.com/Dynatrace/dynatrace-operator/src/dtclient"
	"github.com/Dynatrace/dynatrace-operator/src/kubeobjects"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"golang.org/x/net/http/httpproxy"
	"k8s.io/apimachinery/pkg/types"
)

func checkProxySettings(troubleshootCtx *troubleshoot.TroubleshootContext) error {
	return checkProxySettingsWithLog(troubleshootCtx, troubleshoot.newSubTestLogger("proxy"))
}

func checkProxySettingsWithLog(troubleshootCtx *troubleshoot.TroubleshootContext, logger logr.Logger) error {
	log = logger

	var proxyURL string
	troubleshoot.LogNewCheckf("Analyzing proxy settings ...")

	proxySettingsAvailable := false
	if troubleshootCtx.dynakube.HasProxy() {
		proxySettingsAvailable = true
		troubleshoot.LogInfof("Reminder: Proxy settings in the Dynakube do not apply to pulling of pod images. Please set your proxy on accordingly on node level.")
		troubleshoot.LogWarningf("Proxy settings in the Dynakube are ignored for codeModules images due to technical limitations.")

		var err error
		proxyURL, err = getProxyURL(troubleshootCtx)
		if err != nil {
			troubleshoot.LogErrorf("Unexpected error when reading proxy settings from Dynakube: %v", err)
			return nil
		}
	}

	if checkEnvironmentProxySettings(proxyURL) {
		proxySettingsAvailable = true
	}

	if !proxySettingsAvailable {
		troubleshoot.LogOkf("No proxy settings found.")
	}
	return nil
}

func checkEnvironmentProxySettings(proxyURL string) bool {
	envProxy := getEnvProxySettings()
	if envProxy == nil {
		return false
	}

	troubleshoot.LogInfof("Searching environment for proxy settings ...")
	if envProxy.HTTPProxy != "" {
		troubleshoot.LogWarningf("HTTP_PROXY is set in environment. This setting will be used by the operator for codeModule image pulls.")
		if proxySettingsDiffer(envProxy.HTTPProxy, proxyURL) {
			troubleshoot.LogWarningf("Proxy settings in the Dynakube and HTTP_PROXY differ.")
		}
	}
	if envProxy.HTTPSProxy != "" {
		troubleshoot.LogWarningf("HTTPS_PROXY is set in environment. This setting will be used by the operator for codeModule image pulls.")
		if proxySettingsDiffer(envProxy.HTTPSProxy, proxyURL) {
			troubleshoot.LogWarningf("Proxy settings in the Dynakube and HTTPS_PROXY differ.")
		}
	}
	return true
}

func proxySettingsDiffer(envProxy, dynakubeProxy string) bool {
	return envProxy != "" && dynakubeProxy != "" && envProxy != dynakubeProxy
}

func getEnvProxySettings() *httpproxy.Config {
	envProxy := httpproxy.FromEnvironment()
	if envProxy.HTTPProxy != "" || envProxy.HTTPSProxy != "" {
		return envProxy
	}
	return nil
}

func applyProxySettings(troubleshootCtx *troubleshoot.TroubleshootContext) error {
	proxyURL, err := getProxyURL(troubleshootCtx)
	if err != nil {
		return err
	}

	if proxyURL != "" {
		err := troubleshootCtx.SetTransportProxy(proxyURL)
		if err != nil {
			return errors.Wrapf(err, "error parsing proxy value")
		}
	}

	return nil
}

func getProxyURL(troubleshootCtx *troubleshoot.TroubleshootContext) (string, error) {
	if troubleshootCtx.dynakube.Spec.Proxy == nil {
		return "", nil
	}

	if troubleshootCtx.dynakube.Spec.Proxy.Value != "" {
		return troubleshootCtx.dynakube.Spec.Proxy.Value, nil
	}

	if troubleshootCtx.dynakube.Spec.Proxy.ValueFrom != "" {
		err := setProxySecret(troubleshootCtx)
		if err != nil {
			return "", err
		}

		proxyUrl, err := kubeobjects.ExtractToken(troubleshootCtx.proxySecret, dtclient.CustomProxySecretKey)
		if err != nil {
			return "", errors.Wrapf(err, "failed to extract proxy secret field")
		}
		return proxyUrl, nil
	}
	return "", nil
}

func setProxySecret(troubleshootCtx *troubleshoot.TroubleshootContext) error {
	if troubleshootCtx.proxySecret != nil {
		return nil
	}

	query := kubeobjects.NewSecretQuery(troubleshootCtx.Context, nil, troubleshootCtx.ApiReader, log)
	secret, err := query.Get(types.NamespacedName{
		Namespace: troubleshootCtx.Namespace,
		Name:      troubleshootCtx.dynakube.Spec.Proxy.ValueFrom})

	if err != nil {
		return errors.Wrapf(err, "'%s:%s' proxy secret is missing",
			troubleshootCtx.Namespace, troubleshootCtx.dynakube.Spec.Proxy.ValueFrom)
	}

	troubleshootCtx.proxySecret = &secret
	troubleshoot.LogInfof("proxy secret '%s:%s' exists",
		troubleshootCtx.Namespace, troubleshootCtx.dynakube.Spec.Proxy.ValueFrom)
	return nil
}
