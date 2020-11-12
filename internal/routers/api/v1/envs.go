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

// @Summary 获取helm环境变量
// @Product json
// @Tags	Env
// @Success 200 {object} app.ResponseBody "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/envs [get]
func (e Env) GetHelmEnvs(c *gin.Context) {
	response := app.NewResponse(c)
	response.ToResponse(gin.H{"envs": global.HelmClientSettings.EnvVars()})
	return
}
