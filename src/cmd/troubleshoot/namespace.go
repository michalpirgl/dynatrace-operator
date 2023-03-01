package troubleshoot

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type namespaceCheck struct {
	ctx       context.Context
	namespace string
	apiReader client.Reader
}

func newNamespaceCheck(ctx context.Context, namespace string, apiReader client.Reader) namespaceCheck {
	return namespaceCheck{
		ctx:       ctx,
		namespace: namespace,
		apiReader: apiReader,
	}
}

func (c namespaceCheck) Name() string {
	return "namespaceCheck"
}

func (c namespaceCheck) Do(baseLog logr.Logger) error {
	log := baseLog.WithName(c.Name())

	logNewCheckf(log, "checking if namespace '%s' exists ...", c.namespace)

	var namespace corev1.Namespace
	err := c.apiReader.Get(c.ctx, client.ObjectKey{Name: c.namespace}, &namespace)

	if err != nil {
		return errors.Wrapf(err, "missing namespace '%s'", c.namespace)
	}

	logOkf(log, "using namespace '%s'", c.namespace)
	return nil
}
