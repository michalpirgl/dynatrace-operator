package support_archive

import (
	"context"
	"os"
	"sort"
	"strings"

	"github.com/alecthomas/units"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgocorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const podNameEnvVar = "POD_NAME"
const defaultMemoryLimit = int64(128 * units.MiB)
const operatorContainerName = "dynatrace-operator"

func getMemoryLimit(ctx context.Context, log logr.Logger, pods clientgocorev1.PodInterface) int64 {
	operator, err := getOperatorContainer(ctx, pods)
	if err != nil {
		logErrorf(log, err, "could not read memory limits, using default %d", defaultMemoryLimit)
		return defaultMemoryLimit
	}
	memoryLimits := operator.Resources.Limits.Memory()

	if memoryLimits == nil {
		logErrorf(log, err, "no memory limits defined using default %d", defaultMemoryLimit)
		return defaultMemoryLimit
	}
	return memoryLimits.Value()
}

func getOperatorContainer(ctx context.Context, pods clientgocorev1.PodInterface) (corev1.Container, error) {
	podName := os.Getenv(podNameEnvVar)

	opts := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			Kind: "pod",
		},
		ResourceVersion: "",
	}

	opPod, err := pods.Get(ctx, podName, opts)
	if err != nil {
		return corev1.Container{}, err
	}

	i, found := sort.Find(len(opPod.Spec.Containers), func(i int) int {
		return strings.Compare(operatorContainerName, opPod.Spec.Containers[i].Name)
	})

	if found {
		return opPod.Spec.Containers[i], nil
	}

	return corev1.Container{}, errors.Errorf("Operator container not found in pod %s", podName)
}
