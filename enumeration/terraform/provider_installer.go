package terraform

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	error2 "github.com/snyk/driftctl/enumeration/terraform/error"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type HomeDirInterface interface {
	Dir() (string, error)
}

type ProviderInstaller struct {
	downloader ProviderDownloaderInterface
	config     ProviderConfig
	homeDir    string
}

func NewProviderInstaller(config ProviderConfig) (*ProviderInstaller, error) {
	return &ProviderInstaller{
		NewProviderDownloader(),
		config,
		config.ConfigDir,
	}, nil
}

func (p *ProviderInstaller) Install() (string, error) {
	providerDir := p.getProviderDirectory()
	providerPath := p.getBinaryPath()

	info, err := os.Stat(providerPath)

	if err != nil && os.IsNotExist(err) {
		logrus.WithFields(logrus.Fields{
			"path": providerPath,
		}).Debug("provider not found, downloading ...")
		err := p.downloader.Download(
			p.config.GetDownloadUrl(),
			providerDir,
		)
		if err != nil {
			if notFoundErr, ok := err.(error2.ProviderNotFoundError); ok {
				notFoundErr.Version = p.config.Version
				return "", notFoundErr
			}
			return "", err
		}
		logrus.Debug("Download successful")
	}

	if info != nil && info.IsDir() {
		return "", errors.Errorf(
			"found directory instead of provider binary in %s",
			providerPath,
		)
	}

	if info != nil {
		logrus.WithFields(logrus.Fields{
			"path": providerPath,
		}).Debug("Found existing provider")
	}

	return p.getBinaryPath(), nil
}

func (p ProviderInstaller) getProviderDirectory() string {
	return path.Join(p.homeDir, fmt.Sprintf(".driftctl/plugins/%s_%s/", runtime.GOOS, runtime.GOARCH))
}

// Handle postfixes in binary names
func (p *ProviderInstaller) getBinaryPath() string {
	providerDir := p.getProviderDirectory()
	binaryName := p.config.GetBinaryName()
	_, err := os.Stat(path.Join(providerDir, binaryName))
	if err != nil && os.IsNotExist(err) {
		_ = filepath.WalkDir(providerDir, func(filePath string, d fs.DirEntry, err error) error {
			if d != nil && strings.HasPrefix(d.Name(), p.config.GetBinaryName()) {
				binaryName = d.Name()
			}
			return nil
		})
	}

	return path.Join(providerDir, binaryName)
}
