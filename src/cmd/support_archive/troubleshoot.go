package support_archive

import (
	"context"
	"net/http"

	"github.com/Dynatrace/dynatrace-operator/src/cmd/troubleshoot"
	"github.com/go-logr/logr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const troubleshootCollectorName = "troubleshoot"

type troubleshootCollector struct {
	collectorCommon

	context    context.Context
	apiReader  client.Reader
	kubeConfig rest.Config
	namespace  string
}

func newTroubleshootCollector(context context.Context, log logr.Logger, supportArchive tarball, namespace string, apiReader client.Reader, kubeConfig rest.Config) collector { //nolint:revive // argument-limit doesn't apply to constructors
	return troubleshootCollector{
		collectorCommon: collectorCommon{
			log:            log,
			supportArchive: supportArchive,
		},
		context:    context,
		apiReader:  apiReader,
		kubeConfig: kubeConfig,
		namespace:  namespace,
	}
}

func (t troubleshootCollector) Name() string {
	return troubleshootCollectorName
}

func (t troubleshootCollector) Do() error {

	//TODO: proper constructor for TroubleshootContext
	troubleshootCtx := troubleshoot.TroubleshootContext{
		context:       context.Background(),
		apiReader:     t.apiReader,
		httpClient:    &http.Client{},
		namespaceName: t.namespace,
		kubeConfig:    t.kubeConfig,
	}

	troubleshoot.RunAllChecks(&troubleshootCtx, t.apiReader)

	return nil
}
