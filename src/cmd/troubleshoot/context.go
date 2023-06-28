package troubleshoot

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"
	"github.com/Dynatrace/dynatrace-operator/src/controllers/dynakube/token"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TroubleshootContext struct {
	Context                  context.Context
	ApiReader                client.Reader
	HttpClient               *http.Client
	Namespace                string // the default namespace ("dynatrace") or provided in the command line
	dynakube                 v1beta1.DynaKube
	dynatraceApiSecretTokens token.Tokens
	pullSecret               corev1.Secret
	proxySecret              *corev1.Secret
	KubeConfig               rest.Config
}

func (troubleshootCtx *TroubleshootContext) SetTransportProxy(proxy string) error {
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return errors.Wrap(err, "could not parse proxy URL!")
		}

		if troubleshootCtx.HttpClient.Transport == nil {
			troubleshootCtx.HttpClient.Transport = http.DefaultTransport
		}

		troubleshootCtx.HttpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
		LogInfof("using '%s' proxy to connect to the registry", proxyUrl.Host)
	}
	return nil
}
