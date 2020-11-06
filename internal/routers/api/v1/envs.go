package v1

import (
	"github.com/gin-gonic/gin"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
)

type Env struct {
}

func NewEnv() Env {
	return Env{}
}

func (e Env) GetHelmEnvs(c *gin.Context) {
	app.RespOK(c, global.HelmClientSettings.EnvVars())
}
