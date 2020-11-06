package global

import (
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
	"time"
)

type HelmConfig struct {
	UploadPath   string        `yaml:"uploadPath"`
	HelmRepos    []*repo.Entry `yaml:"helmRepos"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

var (
	HelmClientSettings = cli.New()
	DefaultUploadPath  = "/tmp/charts"
	MyHelmConfig       = &HelmConfig{}
)
