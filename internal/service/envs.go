package service

import (
	"github.com/gin-gonic/gin"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
)

func GetHelmEnvs(c *gin.Context) {
	app.RespOK(c, global.HelmClientSettings.EnvVars())
}
