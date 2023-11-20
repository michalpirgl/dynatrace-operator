//go:build e2e

package injection_failure_policy

import (
	dynatracev1beta1 "github.com/Dynatrace/dynatrace-operator/pkg/api/v1beta1/dynakube"
	corev1 "k8s.io/api/core/v1"
	"testing"

	"github.com/Dynatrace/dynatrace-operator/test/features/cloudnative"
	"github.com/Dynatrace/dynatrace-operator/test/helpers"
	"github.com/Dynatrace/dynatrace-operator/test/helpers/components/dynakube"
	"github.com/Dynatrace/dynatrace-operator/test/helpers/kubeobjects/namespace"
	"github.com/Dynatrace/dynatrace-operator/test/helpers/sample"
	"github.com/Dynatrace/dynatrace-operator/test/helpers/tenant"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

const (
	nodeTaintKey = "injection_failure_policy"
)

func Feature(t *testing.T) features.Feature {
	builder := features.New("cloudnative injection failure policy - force")
	builder.WithLabel("name", "cloudnative-injection-failure-policy-force")

	secretConfig := tenant.GetSingleTenantSecret(t)

	cloudNativeSpec := cloudnative.DefaultCloudNativeSpec()
	cloudNativeSpec.NodeSelector = map[string]string{
		"injection_failure_policy": "e2eaaa",
	}

	cloudNativeSpec.Tolerations = []corev1.Toleration{
		{
			Key:      nodeTaintKey,
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoExecute,
		},
	}

	testDynakube := *dynakube.New(
		dynakube.WithAnnotations(map[string]string{
			//dynatracev1beta1.AnnotationInjectionFailurePolicy: "force",
			dynatracev1beta1.AnnotationInjectionFailurePolicy: "fail",
			//dynatracev1beta1.AnnotationInjectionFailurePolicy: "silent",
		}),
		dynakube.WithApiUrl(secretConfig.ApiUrl),
		dynakube.WithCloudNativeSpec(cloudNativeSpec),
	)

	// Register sample app install
	sampleNamespace := *namespace.New("cloudnative-disabled-injection-sample")
	sampleApp := sample.NewApp(t, &testDynakube,
		sample.AsDeployment(),
		sample.WithNamespace(sampleNamespace),
		sample.WithTolerations([]corev1.Toleration{
			{
				Key:      nodeTaintKey,
				Operator: corev1.TolerationOpExists,
				Effect:   corev1.TaintEffectNoExecute,
			},
		}),
		sample.WithNodeSelector(map[string]string{
			"injection_failure_policy": "e2e",
		}),
	)
	builder.Assess("create sample namespace", sampleApp.InstallNamespace())

	// Register dynakube install
	dynakube.Install(builder, helpers.LevelAssess, &secretConfig, testDynakube)

	// Register sample app install
	builder.Assess("install sample app", sampleApp.Install())

	// Register actual test
	assessSampleInitContainersEnabled(builder, sampleApp)

	// Register sample, dynakube and operator uninstall
	//	builder.Teardown(sampleApp.Uninstall())
	//	dynakube.Delete(builder, helpers.LevelTeardown, testDynakube)
	return builder.Feature()
}

func assessSampleInitContainersEnabled(builder *features.FeatureBuilder, sampleApp *sample.App) {
	builder.Assess("sample apps have init containers", checkInitContainersInjected(sampleApp))
}

func checkInitContainersInjected(sampleApp *sample.App) features.Func {
	return func(ctx context.Context, t *testing.T, envConfig *envconf.Config) context.Context {
		resources := envConfig.Client().Resources()

		pods := sampleApp.GetPods(ctx, t, resources)
		require.NotEmpty(t, pods.Items)

		for _, podItem := range pods.Items {
			if podItem.DeletionTimestamp != nil {
				continue
			}

			require.NotNil(t, podItem)
			require.NotNil(t, podItem.Spec)
			require.NotEmpty(t, podItem.Spec.InitContainers)
		}

		return ctx
	}
}
