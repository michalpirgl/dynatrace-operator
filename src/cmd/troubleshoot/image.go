package troubleshoot

import (
	"fmt"
	"os"
	_ "unsafe"

	"github.com/Dynatrace/dynatrace-operator/src/controllers/dynakube/dtpullsecret"
	"github.com/Dynatrace/dynatrace-operator/src/dockerconfig"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	"golang.org/x/net/http/httpproxy"
)

const (
	pullSecretSuffix = "-pull-secret"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

type Endpoints map[string]Credentials

type Auths struct {
	Auths Endpoints `json:"auths"`
}

func verifyAllImagesAvailable(troubleshootCtx *troubleshootContext) error {
	log := troubleshootCtx.baseLog.WithName("imagepull")

	if troubleshootCtx.dynakube.NeedsOneAgent() {
		verifyImageIsAvailable(log, troubleshootCtx, componentOneAgent, false)
		verifyImageIsAvailable(log, troubleshootCtx, componentCodeModules, true)
	}
	if troubleshootCtx.dynakube.NeedsActiveGate() {
		verifyImageIsAvailable(log, troubleshootCtx, componentActiveGate, false)
	}
	return nil
}

func verifyImageIsAvailable(log logr.Logger, troubleshootCtx *troubleshootContext, comp component, proxyWarning bool) {
	image, isCustomImage := comp.getImage(&troubleshootCtx.dynakube)
	if comp.SkipImageCheck(image) {
		logErrorf(log, "Unknown %s image", comp.String())
		return
	}

	componentName := comp.Name(isCustomImage)
	logNewCheckf(log, "Verifying that %s image %s can be pulled ...", componentName, image)

	if image != "" {
		if troubleshootCtx.dynakube.HasProxy() && proxyWarning {
			logWarningf(log, "Proxy setting in Dynakube is ignored for %s image due to technical limitations.", componentName)
		}

		if getEnvProxySettings() != nil {
			logWarningf(log, "Proxy settings in environment might interfere when pulling %s image in troubleshoot mode.", componentName)
		}

		err := tryImagePull(troubleshootCtx, image)
		if err != nil {
			logErrorf(log, "Pulling %s image %s failed: %v", componentName, image, err)
		} else {
			logOkf(log, "%s image %s can be successfully pulled", componentName, image)
		}
	} else {
		logInfof(log, "No %s image configured", componentName)
	}
}

func tryImagePull(troubleshootCtx *troubleshootContext, image string) error {
	imageReference, err := docker.ParseReference(normalizeDockerReference(image))
	if err != nil {
		return err
	}

	dockerCfg := dockerconfig.NewDockerConfig(troubleshootCtx.apiReader, troubleshootCtx.dynakube)
	defer func(dockerCfg *dockerconfig.DockerConfig, fs afero.Afero) {
		_ = dockerCfg.Cleanup(fs)
	}(dockerCfg, troubleshootCtx.fs)

	systemCtx, err := makeSysContext(troubleshootCtx, imageReference, dockerCfg)
	if err != nil {
		return err
	}
	systemCtx.DockerInsecureSkipTLSVerify = types.OptionalBoolTrue

	restoreProxies := rewireProxy("192.168.0.0")
	imageSource, err := imageReference.NewImageSource(troubleshootCtx.context, systemCtx)
	restoreProxies()

	if err != nil {
		return err
	}

	_ = imageSource.Close()
	return nil
}

//go:linkname resetCachedProxies net/http.resetProxyConfig
func resetCachedProxies()

func rewireProxy(proxy string) func() {
	oldConfig := httpproxy.FromEnvironment()
	oldRequestMethod := os.Getenv("REQUEST_METHOD")

	os.Setenv("HTTP_PROXY", proxy)
	os.Setenv("HTTPS_PROXY", proxy)
	resetCachedProxies()

	return func() {
		os.Setenv("HTTP_PROXY", oldConfig.HTTPProxy)
		os.Setenv("HTTPS_PROXY", oldConfig.HTTPSProxy)
		os.Setenv("NO_PROXY", oldConfig.NoProxy)
		os.Setenv("REQUEST_METHOD", oldRequestMethod)
	}
}

//func IsValueOk(val reflect.Value, expectedType string, expectedKind reflect.Kind) (reflect.Value, error) {
//	if !val.IsValid() {
//		return val, errors.New("value is invalid")
//	}
//	if val.IsZero() {
//		return val, errors.New("value is zero")
//	}
//	if val.IsNil() {
//		return val, errors.New("value is nil")
//	}
//	if val.Type().String() != expectedType {
//		return val, errors.Errorf("expected type %s, got %s", expectedType, val.Type().String())
//	}
//	if val.Kind() != expectedKind {
//		return val, errors.Errorf("expected kind %s, got %s", expectedKind.String(), val.Kind().String())
//	}
//	return val, nil
//}
//
//func injectProxyToImageSource(troubleshootCtx *troubleshootContext, imageSource types.ImageSource) (returnErr error) {
//	defer func() {
//		// we are using reflection and unsafe in this function, better be safe than sorry
//		if err := recover(); err != nil {
//			returnErr = errors.Errorf("caught error while injecting proxy: %v", err)
//		}
//	}()
//
//	//imageSourceVal=&{{0xc0009b38f0} {{true}} {} {} {{{127.0.0.1:46719 linux/oneagent} latest}} {{{127.0.0.1:46719 linux/oneagent} latest}} 0xc00099c780 [] } type=*docker.dockerImageSource kind=ptr
//	//ptrDockerClientVal=&{0xc0009b6000 127.0.0.1:46719 containers/5.25.0 (github.com/containers/image) 0xc0000bd200 {  }  0xc000202510 false {repository linux/oneagent pull} 0xc00047a780 https [] false {{0 0} {[] {} <nil>} map[] 0} {1 {0 0}} <nil>} type=*docker.dockerClient kind=ptr
//	//ptrHttpClientVal=&{0xc00099c8c0 <nil> <nil> 0} type=*http.Client kind=ptr
//	//ptrTransportVal=0xc00099c8c0 type=http.RoundTripper kind=interface
//	//ptrTransportVal=&{{0 0} false map[{ https 127.0.0.1:46719 false}:[0xc00012ab40]] map[] {0xc000334090 map[0xc00012ab40:0xc00047ae40]} {0 0} map[] {0 0} {<nil>} {0 0} map[] map[] 0x1758b20 <nil> 0x7aab60 <nil> <nil> <nil> 0xc0000bd200 10000000000 false false 100 0 0 90000000000 0 0 map[] map[] <nil> 0 0 0 {1 {0 0}} <nil> true false} type=*http.Transport kind=ptr
//
//	proxy, err := getProxyURL(troubleshootCtx)
//	if err != nil {
//		return errors.Wrap(err, "unable to get proxy URL")
//	}
//
//	proxyUrl, err := url.Parse(proxy)
//	if err != nil {
//		return errors.Wrap(err, "unable to get proxy URL")
//	}
//
//	ptrImageSourceVal, err := IsValueOk(reflect.ValueOf(imageSource), "*docker.dockerImageSource", reflect.Pointer)
//	if err != nil {
//		return errors.Wrap(err, "could not reflect on image source")
//	}
//
//	ptrDockerClientVal, err := IsValueOk(ptrImageSourceVal.Elem().FieldByName("c"), "*docker.dockerClient", reflect.Pointer)
//	if err != nil {
//		return errors.Wrap(err, "could not reflect on docker client")
//	}
//
//	ptrHttpClientVal, err := IsValueOk(ptrDockerClientVal.Elem().FieldByName("client"), "*http.Client", reflect.Pointer)
//	if err != nil {
//		return errors.Wrap(err, "could not reflect on http client")
//	}
//
//	ifTransportVal, err := IsValueOk(ptrHttpClientVal.Elem().FieldByName("Transport"), "http.RoundTripper", reflect.Interface)
//	if err != nil {
//		return errors.Wrap(err, "could not reflect on transport interface")
//	}
//
//	ptrTransportVal, err := IsValueOk(ifTransportVal.Elem(), "*http.Transport", reflect.Pointer)
//	if err != nil {
//		return errors.Wrap(err, "could not reflect on transport")
//	}
//
//	// !!!! now it gets sketchy
//	transport := (*http.Transport)(ptrTransportVal.UnsafePointer())
//	transport.Proxy = http.ProxyURL(proxyUrl)
//
//	return nil
//}

func normalizeDockerReference(image string) string {
	return "//" + image
}

func makeSysContext(troubleshootCtx *troubleshootContext, imageReference types.ImageReference, dockerCfg *dockerconfig.DockerConfig) (*types.SystemContext, error) {
	dockerCfg.SetRegistryAuthSecret(&troubleshootCtx.pullSecret)
	err := dockerCfg.StoreRequiredFiles(troubleshootCtx.context, troubleshootCtx.fs)
	if err != nil {
		return nil, err
	}
	return dockerconfig.MakeSystemContext(imageReference.DockerReference(), dockerCfg), nil
}

func getPullSecretToken(troubleshootCtx *troubleshootContext) (string, error) {
	secretBytes, hasPullSecret := troubleshootCtx.pullSecret.Data[dtpullsecret.DockerConfigJson]
	if !hasPullSecret {
		return "", fmt.Errorf("token .dockerconfigjson does not exist in secret '%s'", troubleshootCtx.pullSecret.Name)
	}

	secretStr := string(secretBytes)
	return secretStr, nil
}
