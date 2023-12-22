package image

import (
	"context"
	"fmt"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	"github.com/Dynatrace/dynatrace-operator/pkg/injection/codemodule/installer/common"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	containerv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/pkg/errors"
)

type imagePullInfo struct {
	imageCacheDir string
	targetDir     string
}

func (installer Installer) extractAgentBinariesFromImage(pullInfo imagePullInfo, imageName string) error { //nolint
	//img, err := installer.pullImageInfo(imageName)
	//if err != nil {
	//	log.Info("pullImageInfo", "error", err)
	//	return err
	//}

	//image := *img

	// TODO uncomment => simulates a pull error
	// TODO check if it is a pull secret error and if so, try to pull the image with the kubernetes job method
	//err = installer.pullOCIimage(image, imageName, pullInfo.imageCacheDir, pullInfo.targetDir)
	//if err != nil {
	//	log.Info("pullOCIimage", "err", err)
	//	return err
	//}

	// if the pull fails, we try to let kubernetes itself pull the image and store it in the folder
	return installer.createDownloadJob(imageName, pullInfo.targetDir)
}

func (installer Installer) pullImageInfo(imageName string) (*containerv1.Image, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, errors.WithMessagef(err, "parsing reference %q:", imageName)
	}

	image, err := remote.Image(ref, remote.WithContext(context.TODO()), remote.WithAuthFromKeychain(installer.keychain), remote.WithTransport(installer.transport))
	if err != nil {
		return nil, errors.WithMessagef(err, "getting image %q", imageName)
	}
	return &image, nil
}

func (installer Installer) pullOCIimage(image containerv1.Image, imageName string, imageCacheDir string, targetDir string) error {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return errors.WithMessagef(err, "parsing reference %q", imageName)
	}

	log.Info("pullOciImage", "ref_identifier", ref.Identifier(), "ref.Name", ref.Name(), "ref.String", ref.String())

	err = installer.fs.MkdirAll(imageCacheDir, common.MkDirFileMode)
	if err != nil {
		log.Info("failed to create cache dir", "dir", imageCacheDir, "err", err)
		return errors.WithStack(err)
	}

	if err := crane.SaveOCI(image, path.Join(imageCacheDir, ref.Identifier())); err != nil {
		log.Info("saving v1.Image img as an OCI Image Layout at path", imageCacheDir, err)
		return errors.WithMessagef(err, "saving v1.Image img as an OCI Image Layout at path %s", imageCacheDir)
	}

	layers, err := image.Layers()
	if err != nil {
		log.Info("failed to get image layers", "err", err)
		return errors.WithStack(err)
	}

	err = installer.unpackOciImage(layers, filepath.Join(imageCacheDir, ref.Identifier()), targetDir)
	if err != nil {
		log.Info("failed to unpackOciImage", "error", err)
		return errors.WithStack(err)
	}
	return nil
}

func (installer Installer) createDownloadJob(imageName, imageTargetDir string) error {
	err := installer.fs.MkdirAll(imageTargetDir, common.MkDirFileMode)
	if err != nil {
		log.Info("failed to create cache dir", "dir", imageTargetDir, "err", err)
		return errors.WithStack(err)
	}
	log.Info("createDownloadJob", "imageTargetDir", imageTargetDir)

	backoffLimit := int32(4)
	targetDirVolumeName := "target"

	job := v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      os.Getenv("NODE_NAME"),
			Namespace: os.Getenv("POD_NAMESPACE"),
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": os.Getenv("NODE_NAME"),
					},
					Containers: []corev1.Container{
						{
							Name:  "download-agent",
							Image: imageName,
							//Command: []string{"sleep", "100000"},
							Command: []string{"cp", "-r", "/opt/dynatrace/oneagent/.", "/target"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      targetDirVolumeName,
									MountPath: "/target",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: targetDirVolumeName,
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									// TODO would need to be fetched from csi ds as it can be changed by customer in the helm file
									Path: path.Join("/var/lib/kubelet/plugins/csi.oneagent.dynatrace.com", imageTargetDir),
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}

	err = installer.client.Create(context.Background(), &job)
	if err != nil {
		log.Info("createDownloadJob", "error", err)
		return errors.WithStack(err)
	}
	defer func() {
		_ = installer.client.Delete(context.Background(), &job)
	}()

	// TODO add proper error handling => what if pod can never be pulled? because of missing permissions or not required space for another pod on the node
	// poll every 5 seconds until job is completed
	for true {
		err = installer.client.Get(
			context.Background(),
			client.ObjectKey{Name: job.Name, Namespace: job.Namespace},
			&job,
		)
		if err != nil && !errors2.IsNotFound(err) {
			return errors.WithStack(err)
		}
		if job.Status.Succeeded > 0 {
			log.Info("waiting for download job to complete")
			break
		}
		// wait for 5 seconds
		log.Info("download job is running")
		time.Sleep(10 * time.Second)
	}

	return nil
}

func (installer Installer) unpackOciImage(layers []containerv1.Layer, imageCacheDir string, targetDir string) error {
	for _, layer := range layers {
		mediaType, _ := layer.MediaType()
		switch mediaType {
		case types.DockerLayer:
			digest, _ := layer.Digest()
			sourcePath := filepath.Join(imageCacheDir, "blobs", digest.Algorithm, digest.Hex)
			log.Info("unpackOciImage", "sourcePath", sourcePath)
			if err := installer.extractor.ExtractGzip(sourcePath, targetDir); err != nil {
				return err
			}
		case types.OCILayer:
			return errors.New("OCILayer is not implemented")
		case types.OCILayerZStd:
			return errors.New("OCILayerZStd is not implemented")
		default:
			return fmt.Errorf("media type %s is not implemented", mediaType)
		}
	}
	log.Info("unpackOciImage", "targetDir", targetDir)
	return nil
}
