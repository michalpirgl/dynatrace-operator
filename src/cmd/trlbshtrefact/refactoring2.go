package trlbshtrefact

import (
	"context"
	"log"
)

func checkOneAgentAPM(x ...any) bool {
	return true
}

func checkCRD(x ...any) bool {
	return true
}
func checkNamespace(x ...any) bool {
	return true
}
func checkIfDynatraceApiSecretHasApiToken(x ...any) (any, bool) {
	return nil, true
}
func checkApiUrlSyntax(x ...any) bool {
	return true
}
func checkDynatraceApiTokenScopes(x ...any) bool {
	return true
}
func checkApiUrlForLatestAgentVersion(x ...any) bool {
	return true
}
func checkPullSecretExists(x ...any) (any, bool) {
	return nil, true
}
func checkPullSecretHasRequiredTokens(x ...any) bool {
	return true
}
func verifyAllImagesAvailable(x ...any) bool {
	return true
}
func checkProxySettings(x ...any) bool {
	return true
}

func runChecksFunctional(ctx context.Context, baseLog log.Logger, kubeConfig any, apiReader any, namespace string) {

	oaApmOk := checkOneAgentAPM(baseLog, kubeConfig)
	crdOk := checkCRD(ctx, baseLog, apiReader, namespace)
	nsOk := checkNamespace(baseLog, namespace)

	if !oaApmOk || !crdOk || !nsOk {
		return
	}

	for _, dynakube := range getDynakubesNG() {

		apiUrlOk := checkApiUrlSyntax(dynakube)
		dynatraceApiSecretTokens, apiSecretToknsOk := checkIfDynatraceApiSecretHasApiToken(ctx, baseLog, apiReader, dynakube)

		if !apiUrlOk || !apiSecretToknsOk {
			// skip this dynakube and continue with next one
			continue
		}
		if !checkDynatraceApiTokenScopes(baseLog, apiReader, ctx, dynakube, dynatraceApiSecretTokens) {
			continue
		}
		if !checkApiUrlForLatestAgentVersion(baseLog, ctx, apiReader, dynakube, dynatraceApiSecretTokens) {
			continue
		}
		pullSecret, pullSecretOk := checkPullSecretExists(baseLog, ctx, apiReader, namespace, dynakube)
		if !pullSecretOk {
			continue
		}
		if !checkPullSecretHasRequiredTokens(baseLog, pullSecret) {
			continue
		}

		// ignore errors
		_ = verifyAllImagesAvailable(baseLog, ctx, apiReader, dynakube)
		_ = checkProxySettings(baseLog, dynakube)
	}
}
