package troubleshoot

import (
	"context"
	"fmt"
	"net/http"
	"os"

	dynatracev1beta1 "github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"
	"github.com/Dynatrace/dynatrace-operator/src/cmd/config"
	"github.com/Dynatrace/dynatrace-operator/src/kubeobjects"
	"github.com/Dynatrace/dynatrace-operator/src/scheme"
	"github.com/Dynatrace/dynatrace-operator/src/version"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

const (
	use                    = "troubleshoot"
	dynakubeFlagName       = "dynakube"
	dynakubeFlagShorthand  = "d"
	namespaceFlagName      = "namespace"
	namespaceFlagShorthand = "n"

	namespaceCheckName           = "namespace"
	crdCheckName                 = "crd"
	dynakubeCheckName            = "dynakube"
	oneAgentAPMCheckName         = "oneAgentAPM"
	dtClusterConnectionCheckName = "DynatraceClusterConnection"
	imagePullableCheckName       = "imagePullable"
	proxySettingsCheckName       = "proxySettings"
)

var (
	dynakubeFlagValue  string
	namespaceFlagValue string
)

type CommandBuilder struct {
	configProvider config.Provider
}

func NewTroubleshootCommandBuilder() CommandBuilder {
	return CommandBuilder{}
}

func (builder CommandBuilder) SetConfigProvider(provider config.Provider) CommandBuilder {
	builder.configProvider = provider
	return builder
}

func (builder CommandBuilder) GetCluster(kubeConfig *rest.Config) (cluster.Cluster, error) {
	return cluster.New(kubeConfig, clusterOptions)
}

func (builder CommandBuilder) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:  use,
		RunE: builder.buildRun(),
	}

	addFlags(cmd)

	return cmd
}

func addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&dynakubeFlagValue, dynakubeFlagName, dynakubeFlagShorthand, "", "Specify a different Dynakube name.")
	cmd.PersistentFlags().StringVarP(&namespaceFlagValue, namespaceFlagName, namespaceFlagShorthand, defaultNamespace(), "Specify a different Namespace.")
}

func defaultNamespace() string {
	namespace := os.Getenv("POD_NAMESPACE")

	if namespace == "" {
		return "dynatrace"
	}
	return namespace
}

func clusterOptions(opts *cluster.Options) {
	opts.Scheme = scheme.Scheme
}

func (builder CommandBuilder) buildRun() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		version.LogVersion()

		kubeConfig, err := builder.configProvider.GetConfig()

		if err != nil {
			return err
		}

		err = dynatracev1beta1.AddToScheme(scheme.Scheme)
		if err != nil {
			return err
		}

		k8scluster, err := builder.GetCluster(kubeConfig)
		if err != nil {
			return err
		}

		apiReader := k8scluster.GetAPIReader()

		ctx := context.Background()
		troubleshootCtx := TroubleshootContext{
			Context:    ctx,
			ApiReader:  apiReader,
			HttpClient: &http.Client{},
			Namespace:  namespaceFlagValue,
			KubeConfig: *kubeConfig,
		}

		results := NewChecksResults()
		err = runChecks(results, &troubleshootCtx, getPrerequisiteChecks(ctx, nil, namespaceFlagValue, apiReader)) // ignore error to avoid polluting pretty logs
		resetLogger()
		if err != nil {
			logErrorf("prerequisite checks failed, aborting")
			return err
		}

		dynakubes, err := getDynakubes(context.Background(), apiReader, namespaceFlagValue, dynakubeFlagValue)
		if err != nil {
			return err
		}
		runChecksForAllDynakubes(results, getDynakubeSpecificChecks(results), dynakubes, apiReader)
	}
}

func runChecksForAllDynakubes(results ChecksResults, checks []*CheckListEntry, dynakubes []dynatracev1beta1.DynaKube, apiReader client.Reader) {
	for _, dynakube := range dynakubes {
		results.checkResultMap = map[*CheckListEntry]Result{}
		logNewDynakubef(dynakube.Name)

		troubleshootCtx := TroubleshootContext{
			Context:    context.Background(),
			ApiReader:  apiReader,
			HttpClient: &http.Client{},
			Namespace:  namespaceFlagValue,
			dynakube:   dynakube,
		}

		_ = runChecks(results, &troubleshootCtx, checks) // ignore error to avoid polluting pretty logs
		resetLogger()
		if !results.hasErrors() {
			logOkf("'%s' - all checks passed", dynakube.Name)
		}
	}
}

func getPrerequisiteChecks(ctx context.Context, namespace string, apiReader client.Reader) []*CheckListEntry {
	return []*CheckListEntry{
		{
			Check: newNamespaceCheck(ctx, namespace, apiReader),
		},
		{
			Name: crdCheckName,
			Do:   checkCRD,
		},
		{
			Name: oneAgentAPMCheckName,
			Do:   checkOneAgentAPM,
		},
	}
}

func getDynakubeSpecificChecks(results ChecksResults) []*CheckListEntry {
	checkList := map[string]*CheckListEntry{
		dynakubeCheckName: {
			Do: func(troubleshootCtx *TroubleshootContext) error {
				return checkDynakube(results, troubleshootCtx)
			},
		},
		imagePullableCheckName: {
			Do: verifyAllImagesAvailable,
		},
		proxySettingsCheckName: {
			Do: checkProxySettings,
		},
	}

	checkList[imagePullableCheckName].Prerequisites = []*CheckListEntry{checkList[dynakubeCheckName]}
	checkList[proxySettingsCheckName].Prerequisites = []*CheckListEntry{checkList[dynakubeCheckName]}

	// TODO: return map, entries do not carry their names currently
	return []*CheckListEntry{checkList[dynakubeCheckName], checkList[imagePullableCheckName], checkList[proxySettingsCheckName]}
}

func getDynakubes(ctx context.Context, apiReader client.Reader, namespace string, dynakubeName string) ([]dynatracev1beta1.DynaKube, error) {
	var err error
	var dynakubes []dynatracev1beta1.DynaKube

	if dynakubeName == "" {
		logNewDynakubef("no Dynakube specified - checking all Dynakubes in namespace '%s'", namespace)
		dynakubes, err = getAllDynakubesInNamespace(ctx, apiReader, namespace)
		if err != nil {
			return nil, err
		}
	} else {
		dynakube := dynatracev1beta1.DynaKube{}
		dynakube.Name = dynakubeName
		dynakubes = append(dynakubes, dynakube)
	}

	return dynakubes, nil
}

func getAllDynakubesInNamespace(ctx context.Context, apiReader client.Reader, namespace string) ([]dynatracev1beta1.DynaKube, error) {
	query := kubeobjects.NewDynakubeQuery(apiReader, namespace).WithContext(ctx)
	dynakubes, err := query.List()

	if err != nil {
		logErrorf("failed to list Dynakubes: %v", err)
		return nil, err
	}

	if len(dynakubes.Items) == 0 {
		err = fmt.Errorf("no Dynakubes found in namespace '%s'", namespace)
		logErrorf(err.Error())
		return nil, err
	}

	return dynakubes.Items, nil
}
