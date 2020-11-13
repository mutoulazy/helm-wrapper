package global

import (
	"helm-wrapper/pkg/logger"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
	"time"
)

type HelmConfig struct {
	UploadPath            string        `yaml:"uploadPath"`
	HelmRepos             []*repo.Entry `yaml:"helmRepos"`
	ReadTimeout           time.Duration `yaml:"readTimeout"`
	WriteTimeout          time.Duration `yaml:"writeTimeout"`
	LogSavePath           string        `yaml:"logSavePath"`
	LogFileName           string        `yaml:"logFileName"`
	LogFileExt            string        `yaml:"logFileExt"`
	DefaultContextTimeout time.Duration `yaml:"defaultContextTimeout"`
	RunMode               string        `yaml:"runMode"`
}

var (
	HelmClientSettings = cli.New()
	DefaultUploadPath  = "/tmp/charts"
	MyHelmConfig       = &HelmConfig{}
	Logger             *logger.Logger
)
