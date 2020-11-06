package service

import (
	"helm-wrapper/global"
	"os"

	"github.com/golang/glog"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
)

type KubeInformation struct {
	AimNamespace string
	AimContext   string
}

func InitKubeInformation(namespace, context string) *KubeInformation {
	return &KubeInformation{
		AimNamespace: namespace,
		AimContext:   context,
	}
}

func ActionConfigInit(kubeInfo *KubeInformation) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	if kubeInfo.AimContext == "" {
		kubeInfo.AimContext = global.HelmClientSettings.KubeContext
	}
	clientConfig := kube.GetConfig(global.HelmClientSettings.KubeConfig, kubeInfo.AimContext, kubeInfo.AimNamespace)
	if global.HelmClientSettings.KubeToken != "" {
		clientConfig.BearerToken = &global.HelmClientSettings.KubeToken
	}
	if global.HelmClientSettings.KubeAPIServer != "" {
		clientConfig.APIServer = &global.HelmClientSettings.KubeAPIServer
	}
	err := actionConfig.Init(clientConfig, kubeInfo.AimNamespace, os.Getenv("HELM_DRIVER"), glog.Infof)
	if err != nil {
		glog.Errorf("%+v", err)
		return nil, err
	}

	return actionConfig, nil
}
