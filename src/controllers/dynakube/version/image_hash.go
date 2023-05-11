package version

import (
	"context"
	"fmt"

	"github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"
	"github.com/Dynatrace/dynatrace-operator/src/dockerconfig"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// VersionLabel is the name of the label used on ActiveGate-provided images.
	VersionLabel = "com.dynatrace.build-version"
)

type ImageVersion struct {
	Version string
	Digest  digest.Digest
}

// ImageVersionFunc can fetch image information from img
type ImageVersionFunc func(ctx context.Context, imageName string, dockerConfig *dockerconfig.DockerConfig, kubeReader client.Reader, dynakube v1beta1.DynaKube) (ImageVersion, error)

var _ ImageVersionFunc = GetImageVersion

// GetImageVersion fetches image information for imageName
func GetImageVersion(ctx context.Context, imageName string, dockerConfig *dockerconfig.DockerConfig, kubeReader client.Reader, dynakube v1beta1.DynaKube) (ImageVersion, error) {

	if dynakube.HasProxy() {
		proxy, err := dynakube.Proxy(ctx, kubeReader)
		if err != nil {
			return ImageVersion{}, errors.WithStack(err)
		}
		restoreProxy := OverrideProxyInEnvironment(proxy)
		defer restoreProxy()
	}

	transportImageName := fmt.Sprintf("docker://%s", imageName)

	imageReference, err := alltransports.ParseImageName(transportImageName)
	if err != nil {
		return ImageVersion{}, errors.WithStack(err)
	}

	systemContext := dockerconfig.MakeSystemContext(imageReference.DockerReference(), dockerConfig)

	imageSource, err := imageReference.NewImageSource(ctx, systemContext)
	if err != nil {
		return ImageVersion{}, errors.WithStack(err)
	}
	defer closeImageSource(imageSource)

	imageManifest, _, err := imageSource.GetManifest(ctx, nil)
	if err != nil {
		return ImageVersion{}, errors.WithStack(err)
	}

	digest, err := manifest.Digest(imageManifest)
	if err != nil {
		return ImageVersion{}, errors.WithStack(err)
	}

	sourceImage, err := image.FromUnparsedImage(ctx, systemContext, image.UnparsedInstance(imageSource, nil))
	if err != nil {
		return ImageVersion{}, errors.WithStack(err)
	}

	inspectedImage, err := sourceImage.Inspect(ctx)
	if err != nil {
		return ImageVersion{}, errors.WithStack(err)
	} else if inspectedImage == nil {
		return ImageVersion{}, errors.Errorf("could not inspect image: '%s'", transportImageName)
	}

	return ImageVersion{
		Digest:  digest,
		Version: inspectedImage.Labels[VersionLabel], // empty if unset
	}, nil
}

func closeImageSource(source types.ImageSource) {
	if source != nil {
		// Swallow error
		_ = source.Close()
	}
}
