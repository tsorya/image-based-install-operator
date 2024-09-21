package ops

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	apicfgv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/image-based-install-operator/internal/filelock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	templateExtract = "oc adm release ops --command=%s --to=%s --insecure=%t --registry-config=%s %s"
	installer       = "openshift-install"
)

type Ops struct {
	executer Executer
}

func NewOps(executer Executer) *Ops {
	return &Ops{executer: executer}
}

func getReleaseBinaryPath(releaseImage string, cacheDir string) (workdir string, path string, err error) {
	workdir = filepath.Join(cacheDir, releaseImage)
	path = filepath.Join(workdir, installer)
	return
}

func (o *Ops) Extract(
	log logrus.FieldLogger,
	releaseImage, dir, pullSecret, caBundle string,
	imageDigestMirrors []apicfgv1.ImageDigestMirrors) (string, error) {

	var path string
	var err error
	if releaseImage == "" {
		return "", errors.New("no releaseImage provided")
	}
	registryFile := ""
	if len(imageDigestMirrors) > 0 {
		registryFile, err = o.writeImageDigestSourceToFile(log, imageDigestMirrors)
		if err != nil {
			return "", errors.Wrap(err, "failed to create file ICSP file from registries config")
		}
	}

	caBundleFile, err := o.writeCaBundleToFile(log, caBundle)
	if err != nil {
		return "", errors.Wrap(err, "failed to create ca-bundle file")
	}

	// TODO: check how to handle ca bundle for mirror registries
	path, err = o.extractFromRelease(log, releaseImage, dir, pullSecret, true, registryFile, caBundleFile)
	if err != nil {
		log.WithError(err).Errorf("failed to ops openshift-baremetal-install from release image %s", releaseImage)
		return "", err
	}

	return path, err
}

func (o *Ops) CreateConfigurationIso(log logrus.FieldLogger, installerPath string, workdir string) error {
	log.Infof("Running openshift-install with workdir %s", workdir)
	_, stderr, exitCode := o.executer.Execute(installerPath, []string{"image-based", "create", "config-image", "--dir", workdir}...)
	if exitCode != 0 {
		return fmt.Errorf("failed to create configuration iso: %s", stderr)
	}
	return nil
}
func (o *Ops) extractFromRelease(log logrus.FieldLogger, releaseImage, cacheDir, pullSecret string, insecure bool, registryFile string, caBundleFile string) (string, error) {
	workdir, path, err := getReleaseBinaryPath(releaseImage, cacheDir)
	if err != nil {
		return "", err
	}
	log.Infof("extracting %s binary to %s", path, workdir)
	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path); err == nil {
		log.Infof("binary already exists at %s", path)
		return path, nil
	}
	locked, lockErr, funcErr := filelock.WithWriteLock(workdir, func() error {
		if _, err := os.Stat(path); err == nil {
			log.Infof("binary already exists at %s", path)
			return nil
		}

		var cmd string
		extractCommand := templateExtract
		if registryFile != "" {
			extractCommand = fmt.Sprintf("%s --idms-file='%s'", extractCommand, registryFile)
			defer os.Remove(registryFile)
		}

		if caBundleFile != "" {
			extractCommand = fmt.Sprintf("%s --certificate-authority='%s'", extractCommand, caBundleFile)
			defer os.Remove(caBundleFile)
		}

		pullSecretFile, err := o.createPullSecretFiles(pullSecret)
		if err != nil {
			return fmt.Errorf("failed to create pull-secret file: %w", err)
		}

		cmd = fmt.Sprintf(extractCommand, installer, workdir, insecure, pullSecretFile, releaseImage)
		args := strings.Split(cmd, " ")
		stdout, stderr, exitCode := o.executer.Execute(args[0], args[1:]...)
		if exitCode != 0 {
			err = fmt.Errorf("command '%s' exited with non-zero exit code %d: %s\n%s", cmd, exitCode, stdout, stderr)
			log.Warn(err)
			return err
		}

		return nil
	})

	if lockErr != nil {
		return "", fmt.Errorf("failed to acquire file lock: %w", lockErr)
	}
	if funcErr != nil {
		return "", fmt.Errorf("failed to ops %s input data: %w", installer, funcErr)
	}
	if !locked {
		return "", nil
	}

	log.Infof("Successfully extracted %s binary from the release to: %s", installer, path)
	return path, nil
}

func (o *Ops) createPullSecretFiles(pullSecret string) (string, error) {
	return CreateTempFile(pullSecret, "pull-secret")
}

func (o *Ops) writeImageDigestSourceToFile(log logrus.FieldLogger, imageDigestMirrors []apicfgv1.ImageDigestMirrors) (string, error) {
	if imageDigestMirrors == nil {
		return "", nil
	}

	log.Debugf("Building image-digest-mirror file")

	imageDigestMirrorSet := &apicfgv1.ImageDigestMirrorSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apicfgv1.GroupVersion.String(),
			Kind:       "ImageDigestMirrorSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "image-digest-mirror",
			// not namespaced
		},
		Spec: apicfgv1.ImageDigestMirrorSetSpec{
			ImageDigestMirrors: imageDigestMirrors,
		},
	}

	data, err := json.Marshal(imageDigestMirrorSet)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ImageDigestMirrorSet: %w", err)
	}

	return CreateTempFile(string(data), "image-digest-mirror")
}

func (o *Ops) writeCaBundleToFile(log logrus.FieldLogger, caBundle string) (string, error) {
	if caBundle == "" {
		return "", nil
	}
	log.Debugf("Creating ca-bundle file")

	return CreateTempFile(caBundle, "ca-bundle")
}

func CreateTempFile(data, filePrefix string) (string, error) {
	tmpFile, err := os.CreateTemp("", filePrefix)
	defer tmpFile.Close()
	if err != nil {
		return "", err
	}

	if _, err := tmpFile.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("failed to write %s file: %w", filePrefix, err)
	}

	return tmpFile.Name(), nil
}
