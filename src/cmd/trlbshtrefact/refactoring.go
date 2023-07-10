package trlbshtrefact

import (
	"context"
	"log"
)

type Checklist struct {
}

func (c *Checklist) add(x ...any) {
}

func (c *Checklist) runTests(x ...any) {
}

func newOneAgentAPMCheck(x ...any) interface{} {
	return nil
}

func newCRDCheck(x ...any) interface{} {
	return nil
}
func newNamespaceCheck(x ...any) interface{} {
	return nil
}
func newDtASTC(x ...any) interface{} {
	return nil
}
func newApiUrlSyntaxCheck(x ...any) interface{} {
	return nil
}
func newApiTokenScopeCheck(x ...any) interface{} {
	return nil
}
func newApiUrlLatestAgentVersionCheck(x ...any) interface{} {
	return nil
}
func newPullSecretsCheck(x ...any) interface{} {
	return nil
}
func newPullSecretTokenCheck(x ...any) interface{} {
	return nil
}
func newImagePullableCheck(x ...any) interface{} {
	return nil
}
func newProxySettingsCheck(x ...any) interface{} {
	return nil
}

func getDynakubesNG() []any {
	return nil
}

type check interface {
	Name() string
	Do(log log.Logger) bool
}

func runChecksNG(ctx context.Context, baseLog log.Logger, kubeConfig any, apiReader any, namespace string) {
	checkList := Checklist{}
	oneAgentAPMCheck := newOneAgentAPMCheck(kubeConfig)
	crdCheck := newCRDCheck(ctx, apiReader, namespace)
	namespaceCheck := newNamespaceCheck(namespace)

	checkList.add(oneAgentAPMCheck)
	checkList.add(crdCheck)
	checkList.add(namespaceCheck)

	for _, dynakube := range getDynakubesNG() {
		//  dynakubeCheck := newDynakube(ctx, apiReader, namespace, dynakube, dk)

		dynatraceApiSecretTokenCheck := newDtASTC(ctx, apiReader, dynakube)
		apiUrlSyntaxCheck := newApiUrlSyntaxCheck(dynakube)

		apiTokenScopeCheck := newApiTokenScopeCheck(ctx, apiReader, dynakube, dynatraceApiSecretTokenCheck)
		apiUrlForLatestAgentVersionCheck := newApiUrlLatestAgentVersionCheck(ctx, apiReader, dynakube, dynatraceApiSecretTokenCheck)
		pullSecretCheck := newPullSecretsCheck(ctx, apiReader, namespace, dynakube)
		pullSecretTokenCheck := newPullSecretTokenCheck(pullSecretCheck)

		imagePullableCheck := newImagePullableCheck(ctx, apiReader, dynakube)
		proxySettingsCheck := newProxySettingsCheck(dynakube)

		checkList.add(dynatraceApiSecretTokenCheck, oneAgentAPMCheck, crdCheck, namespace)
		checkList.add(apiUrlSyntaxCheck, oneAgentAPMCheck, crdCheck, namespace)
		checkList.add(apiTokenScopeCheck, dynatraceApiSecretTokenCheck, apiUrlSyntaxCheck)
		checkList.add(apiUrlForLatestAgentVersionCheck, apiTokenScopeCheck)
		checkList.add(pullSecretCheck, apiUrlForLatestAgentVersionCheck)
		checkList.add(pullSecretTokenCheck, pullSecretCheck)
		checkList.add(imagePullableCheck, pullSecretTokenCheck)
		checkList.add(proxySettingsCheck, pullSecretTokenCheck)
	}

	checkList.runTests(baseLog)
}
