package troubleshoot

import (
	"context"
	"fmt"

	"github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"
	"github.com/Dynatrace/dynatrace-operator/src/controllers/dynakube/dtpullsecret"
	"github.com/Dynatrace/dynatrace-operator/src/controllers/dynakube/dynatraceclient"
	"github.com/Dynatrace/dynatrace-operator/src/controllers/dynakube/token"
	"github.com/Dynatrace/dynatrace-operator/src/dtclient"
	"github.com/Dynatrace/dynatrace-operator/src/kubeobjects"
	"github.com/Dynatrace/dynatrace-operator/src/webhook/validation"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	pullSecretFieldValue = "top-secret"

	getSelectedDynakubeCheckName           = "getSelectedDynakube"
	apiUrlSyntaxCheckName                  = "apiUrlSyntax"
	dynatraceApiTokenScopesCheckName       = "dynatraceApiTokenScopes"
	apiUrlLatestAgentVersionCheckName      = "apiUrlLatestAgentVersion"
	dynatraceApiSecretHasApiTokenCheckName = "dynatraceApiSecretHasApiToken"
	pullSecretExistsCheckName              = "pullSecretExists"
	pullSecretHasRequiredTokensCheckName   = "pullSecretHasRequiredTokens"
	proxySecretCheckName                   = "proxySecret"
)

type dynakubeCheck struct {
	ctx       context.Context
	apiReader client.Reader
	namespace string
	dynaKube  v1beta1.DynaKube
}

func newDynakubeCheck(ctx context.Context, apiReader client.Reader, namespace string, dynaKube v1beta1.DynaKube) dynakubeCheck {
	return dynakubeCheck{
		ctx:       ctx,
		apiReader: apiReader,
		namespace: namespace,
		dynaKube:  dynaKube,
	}
}

func (c dynakubeCheck) Name() string {
	return "dynakube"
}

func (c dynakubeCheck) Do() error {

}

func checkDynakube(results ChecksResults, troubleshootCtx *TroubleshootContext) error {
	log = newSubTestLogger("dynakube")

	logNewCheckf("checking if '%s:%s' Dynakube is configured correctly", troubleshootCtx.Namespace, troubleshootCtx.dynakube.Name)

	err := runChecks(results, troubleshootCtx, getDynakubeChecks())
	if err != nil {
		return errors.Wrapf(err, "'%s:%s' Dynakube isn't valid. %s",
			troubleshootCtx.Namespace, troubleshootCtx.dynakube.Name, dynakubeNotValidMessage())
	}

	logOkf("'%s:%s' Dynakube is valid", troubleshootCtx.Namespace, troubleshootCtx.dynakube.Name)
	return nil
}

func getDynakubeChecks() []*CheckListEntry {
	selectedDynakubeCheck := &CheckListEntry{
		Name: getSelectedDynakubeCheckName,
		Do:   getSelectedDynakube,
	}

	ifDynatraceApiSecretHasApiTokenCheck := &CheckListEntry{
		Name:          dynatraceApiSecretHasApiTokenCheckName,
		Do:            checkIfDynatraceApiSecretHasApiToken,
		Prerequisites: []*CheckListEntry{selectedDynakubeCheck},
	}

	apiUrlSyntaxCheck := &CheckListEntry{
		Name:          apiUrlSyntaxCheckName,
		Do:            checkApiUrlSyntax,
		Prerequisites: []*CheckListEntry{selectedDynakubeCheck},
	}

	apiUrlTokenScopesCheck := &CheckListEntry{
		Name:          dynatraceApiTokenScopesCheckName,
		Do:            checkDynatraceApiTokenScopes,
		Prerequisites: []*CheckListEntry{apiUrlSyntaxCheck, ifDynatraceApiSecretHasApiTokenCheck},
	}

	apiUrlLatestAgentVersionCheck := &CheckListEntry{
		Name:          apiUrlLatestAgentVersionCheckName,
		Do:            checkApiUrlForLatestAgentVersion,
		Prerequisites: []*CheckListEntry{apiUrlTokenScopesCheck},
	}

	pullSecretExistsCheck := &CheckListEntry{
		Name:          pullSecretExistsCheckName,
		Do:            checkPullSecretExists,
		Prerequisites: []*CheckListEntry{apiUrlLatestAgentVersionCheck},
	}

	pullSecretHasRequiredTokensCheck := &CheckListEntry{
		Name:          pullSecretHasRequiredTokensCheckName,
		Do:            checkPullSecretHasRequiredTokens,
		Prerequisites: []*CheckListEntry{pullSecretExistsCheck},
	}

	proxySecretIfItExistsCheck := &CheckListEntry{
		Name:          proxySecretCheckName,
		Do:            applyProxySettings,
		Prerequisites: []*CheckListEntry{selectedDynakubeCheck},
	}

	return []*CheckListEntry{selectedDynakubeCheck, ifDynatraceApiSecretHasApiTokenCheck, apiUrlSyntaxCheck, apiUrlTokenScopesCheck, apiUrlLatestAgentVersionCheck, pullSecretExistsCheck, pullSecretHasRequiredTokensCheck, proxySecretIfItExistsCheck}
}

func dynakubeNotValidMessage() string {
	return fmt.Sprintf(
		"Target namespace and dynakube can be changed by providing '--%s <namespace>' or '--%s <dynakube>' parameters.",
		namespaceFlagName, dynakubeFlagName)
}

func getSelectedDynakube(troubleshootCtx *TroubleshootContext) error {
	query := kubeobjects.NewDynakubeQuery(troubleshootCtx.ApiReader, troubleshootCtx.Namespace).WithContext(troubleshootCtx.Context)
	dynakube, err := query.Get(types.NamespacedName{Namespace: troubleshootCtx.Namespace, Name: troubleshootCtx.dynakube.Name})

	if err != nil {
		return determineSelectedDynakubeError(troubleshootCtx, err)
	}

	troubleshootCtx.dynakube = dynakube

	logInfof("using '%s:%s' Dynakube", troubleshootCtx.Namespace, troubleshootCtx.dynakube.Name)
	return nil
}

func determineSelectedDynakubeError(troubleshootCtx *TroubleshootContext, err error) error {
	if k8serrors.IsNotFound(err) {
		err = errors.Wrapf(err,
			"selected '%s:%s' Dynakube does not exist",
			troubleshootCtx.Namespace, troubleshootCtx.dynakube.Name)
	} else {
		err = errors.Wrapf(err, "could not get Dynakube '%s:%s'",
			troubleshootCtx.Namespace, troubleshootCtx.dynakube.Name)
	}
	return err
}

func checkIfDynatraceApiSecretHasApiToken(troubleshootCtx *TroubleshootContext) error {
	tokenReader := token.NewReader(troubleshootCtx.ApiReader, &troubleshootCtx.dynakube)
	tokens, err := tokenReader.ReadTokens(troubleshootCtx.Context)
	if err != nil {
		return errors.Wrapf(err, "'%s:%s' secret is missing or invalid", troubleshootCtx.Namespace, troubleshootCtx.dynakube.Tokens())
	}

	_, hasApiToken := tokens[dtclient.DynatraceApiToken]
	if !hasApiToken {
		return errors.New(fmt.Sprintf("'%s' token is missing in '%s:%s' secret", dtclient.DynatraceApiToken, troubleshootCtx.Namespace, troubleshootCtx.dynakube.Tokens()))
	}

	troubleshootCtx.dynatraceApiSecretTokens = tokens

	logInfof("secret token 'apiToken' exists")
	return nil
}

func checkApiUrlSyntax(troubleshootCtx *TroubleshootContext) error {
	logInfof("checking if syntax of API URL is valid")

	validation.SetLogger(log)
	if validation.NoApiUrl(nil, &troubleshootCtx.dynakube) != "" {
		return errors.New("API URL is invalid")
	}
	if validation.IsInvalidApiUrl(nil, &troubleshootCtx.dynakube) != "" {
		return errors.New("API URL is invalid")
	}

	logInfof("syntax of API URL is valid")
	return nil
}

func checkDynatraceApiTokenScopes(troubleshootCtx *TroubleshootContext) error {
	logInfof("checking if token scopes are valid")

	dtc, err := dynatraceclient.NewBuilder(troubleshootCtx.ApiReader).
		SetContext(troubleshootCtx.Context).
		SetDynakube(troubleshootCtx.dynakube).
		SetTokens(troubleshootCtx.dynatraceApiSecretTokens).
		Build()

	if err != nil {
		return errors.Wrap(err, "failed to build DynatraceAPI client")
	}

	tokens := troubleshootCtx.dynatraceApiSecretTokens.SetScopesForDynakube(troubleshootCtx.dynakube)

	if err = tokens.VerifyValues(); err != nil {
		return errors.Wrapf(err, "invalid '%s:%s' secret", troubleshootCtx.Namespace, troubleshootCtx.dynakube.Tokens())
	}

	if err = tokens.VerifyScopes(dtc); err != nil {
		return errors.Wrapf(err, "invalid '%s:%s' secret", troubleshootCtx.Namespace, troubleshootCtx.dynakube.Tokens())
	}

	logInfof("token scopes are valid")
	return nil
}

func checkApiUrlForLatestAgentVersion(troubleshootCtx *TroubleshootContext) error {
	logInfof("checking if can pull latest agent version")

	dtc, err := dynatraceclient.NewBuilder(troubleshootCtx.ApiReader).
		SetContext(troubleshootCtx.Context).
		SetDynakube(troubleshootCtx.dynakube).
		SetTokens(troubleshootCtx.dynatraceApiSecretTokens).
		Build()
	if err != nil {
		return errors.Wrap(err, "failed to build DynatraceAPI client")
	}

	_, err = dtc.GetLatestAgentVersion(dtclient.OsUnix, dtclient.InstallerTypeDefault)
	if err != nil {
		return errors.Wrap(err, "failed to connect to DynatraceAPI")
	}

	logInfof("API token is valid, can pull latest agent version")
	return nil
}

func checkPullSecretExists(troubleshootCtx *TroubleshootContext) error {
	query := kubeobjects.NewSecretQuery(troubleshootCtx.Context, nil, troubleshootCtx.ApiReader, log)
	secret, err := query.Get(types.NamespacedName{Namespace: troubleshootCtx.Namespace, Name: troubleshootCtx.dynakube.PullSecret()})

	if err != nil {
		return errors.Wrapf(err, "'%s:%s' pull secret is missing", troubleshootCtx.Namespace, troubleshootCtx.dynakube.PullSecret())
	} else {
		troubleshootCtx.pullSecret = secret
	}

	logInfof("pull secret '%s:%s' exists", troubleshootCtx.Namespace, troubleshootCtx.dynakube.PullSecret())
	return nil
}

func checkPullSecretHasRequiredTokens(troubleshootCtx *TroubleshootContext) error {
	if _, err := kubeobjects.ExtractToken(&troubleshootCtx.pullSecret, dtpullsecret.DockerConfigJson); err != nil {
		return errors.Wrapf(err, "invalid '%s:%s' secret", troubleshootCtx.Namespace, troubleshootCtx.dynakube.PullSecret())
	}

	logInfof("secret token '%s' exists", dtpullsecret.DockerConfigJson)
	return nil
}
