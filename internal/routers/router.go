package routers

import (
	"github.com/gin-gonic/gin"
	"helm-wrapper/internal/routers/api/v1"
)

func RegisterRouter(router *gin.Engine) {
	env := v1.NewEnv()
	repository := v1.NewRepository()
	chart := v1.NewChart()
	release := v1.NewRelease()
	apiv1 := router.Group("/api/v1")
	{
		apiv1.GET("/envs", env.GetHelmEnvs)

		apiv1.GET("/repositories/charts", repository.ListRepoCharts)
		apiv1.PUT("/repositories", repository.UpdateRepositories)

		apiv1.GET("/charts", chart.ShowChartInfo)
		apiv1.POST("/charts/upload", UploadChart)
		apiv1.GET("/charts/upload", ListUploadedCharts)

		apiv1.GET("/namespaces/:namespace/releases", release.ListReleases)
		apiv1.GET("/namespaces/:namespace/releases/:release", release.ShowReleaseInfo)
		apiv1.POST("/namespaces/:namespace/releases/:release", release.InstallRelease)
		apiv1.PUT("/namespaces/:namespace/releases/:release", release.UpgradeRelease)
		apiv1.DELETE("/namespaces/:namespace/releases/:release", release.UninstallRelease)
		apiv1.PUT("/namespaces/:namespace/releases/:release/versions/:reversion", release.RollbackRelease)
		apiv1.GET("/namespaces/:namespace/releases/:release/status", release.GetReleaseStatus)
		apiv1.GET("/namespaces/:namespace/releases/:release/histories", release.ListReleaseHistories)
	}
}
