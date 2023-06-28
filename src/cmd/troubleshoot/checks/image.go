package checks

import (
	"context"
	"fmt"
	"github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"
	"github.com/Dynatrace/dynatrace-operator/src/cmd/troubleshoot"
	"github.com/Dynatrace/dynatrace-operator/src/controllers/dynakube/dtpullsecret"
	"github.com/Dynatrace/dynatrace-operator/src/dockerconfig"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	pullSecretSuffix = "-pull-secret"
)

type imageCheck struct {
	ctx        context.Context
	dynaKube   v1beta1.DynaKube
	pullSecret corev1.Secret
	apiReader  client.Reader

	log logr.Logger
}

func newImageCheck(ctx context.Context, dynaKube v1beta1.DynaKube, pullSecret corev1.Secret, apiReader client.Reader) imageCheck {
	return imageCheck{
		ctx:        ctx,
		dynaKube:   dynaKube,
		pullSecret: pullSecret,
		apiReader:  apiReader,
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

type Endpoints map[string]Credentials

type Auths struct {
	Auths Endpoints `json:"auths"`
}

func (c imageCheck) Name() string {
	return "imageCheck"
}

func (c imageCheck) Do(baseLog logr.Logger) error {
	c.log = baseLog.WithName(c.Name())

	if c.dynaKube.NeedsOneAgent() {
		c.verifyImageIsAvailable(troubleshoot.ComponentOneAgent, false)
		c.verifyImageIsAvailable(troubleshoot.ComponentCodeModules, true)
	}
	if c.dynaKube.NeedsActiveGate() {
		c.verifyImageIsAvailable(troubleshoot.ComponentActiveGate, false)
	}
	return nil
}

func (c imageCheck) verifyImageIsAvailable(comp troubleshoot.Component, proxyWarning bool) {
	image, isCustomImage := comp.GetImage(&c.dynaKube)
	if comp.SkipImageCheck(image) {
		troubleshoot.LogErrorf(c.log, "Unknown %s image", comp.String())
		return
	}

	componentName := comp.Name(isCustomImage)
	troubleshoot.LogNewCheckf(c.log, "Verifying that %s image %s can be pulled ...", componentName, image)

	if image != "" {
		if c.dynaKube.HasProxy() && proxyWarning {
			troubleshoot.LogWarningf(c.log, "Proxy setting in Dynakube is ignored for %s image due to technical limitations.", componentName)
		}

		if getEnvProxySettings() != nil {
			troubleshoot.LogWarningf(c.log, "Proxy settings in environment might interfere when pulling %s image in troubleshoot mode.", componentName)
		}

		err := c.tryImagePull(image)
		if err != nil {
			troubleshoot.LogErrorf(c.log, "Pulling %s image %s failed: %v", componentName, image, err)
		} else {
			troubleshoot.LogOkf(c.log, "%s image %s can be successfully pulled", componentName, image)
		}
	} else {
		troubleshoot.LogInfof(c.log, "No %s image configured", componentName)
	}
}

func (c imageCheck) tryImagePull(image string) error {
	imageReference, err := docker.ParseReference(normalizeDockerReference(image))
	if err != nil {
		return err
	}

	systemCtx, err := c.makeSysContext(imageReference)
	if err != nil {
		return err
	}
	systemCtx.DockerInsecureSkipTLSVerify = types.OptionalBoolTrue

	imageSource, err := imageReference.NewImageSource(c.ctx, systemCtx)
	if err != nil {
		return err
	}
	defer imageSource.Close()

	return nil
}

func (c imageCheck) makeSysContext(imageReference types.ImageReference) (*types.SystemContext, error) {
	dockerCfg := dockerconfig.NewDockerConfig(c.apiReader, c.dynaKube)
	err := dockerCfg.SetupAuthsFromSecret(&c.pullSecret)
	if err != nil {
		return nil, err
	}
	return dockerconfig.MakeSystemContext(imageReference.DockerReference(), dockerCfg), nil
}

func (c imageCheck) getPullSecretToken(troubleshootCtx *troubleshoot.TroubleshootContext) (string, error) {
	secretBytes, hasPullSecret := c.pullSecret.Data[dtpullsecret.DockerConfigJson]
	if !hasPullSecret {
		return "", fmt.Errorf("token .dockerconfigjson does not exist in secret '%s'", c.pullSecret.Name)
	}

	secretStr := string(secretBytes)
	return secretStr, nil
}

func normalizeDockerReference(image string) string {
	return "//" + image
}
