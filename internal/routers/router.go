package routers

import (
	"github.com/gin-gonic/gin"
	"helm-wrapper/internal/service"
	"helm-wrapper/pkg/upload"
)

func RegisterRouter(router *gin.Engine) {
	// helm env
	envs := router.Group("/api/envs")
	{
		envs.GET("", service.GetHelmEnvs)
	}

	// helm repo
	repositories := router.Group("/api/repositories")
	{
		// helm search repo
		repositories.GET("/charts", service.ListRepoCharts)
		// helm repo update
		repositories.PUT("", service.UpdateRepositories)
	}

	// helm chart
	charts := router.Group("/api/charts")
	{
		// helm show
		charts.GET("", service.ShowChartInfo)
		// upload chart
		charts.POST("/upload", upload.UploadChart)
		// list uploaded charts
		charts.GET("/upload", upload.ListUploadedCharts)
	}

	// helm release
	releases := router.Group("/api/namespaces/:namespace/releases")
	{
		// helm list releases ->  helm list
		releases.GET("", service.ListReleases)
		// helm get
		releases.GET("/:release", service.ShowReleaseInfo)
		// helm install
		releases.POST("/:release", service.InstallRelease)
		// helm upgrade
		releases.PUT("/:release", service.UpgradeRelease)
		// helm uninstall
		releases.DELETE("/:release", service.UninstallRelease)
		// helm rollback
		releases.PUT("/:release/versions/:reversion", service.RollbackRelease)
		// helm status <RELEASE_NAME>
		releases.GET("/:release/status", service.GetReleaseStatus)
		// helm release history
		releases.GET("/:release/histories", service.ListReleaseHistories)
	}
}
