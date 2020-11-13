package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "helm-wrapper/docs"
	"helm-wrapper/global"
	"helm-wrapper/internal/middleware"
	"helm-wrapper/internal/routers/api/v1"
	"time"
)

func RegisterRouter(router *gin.Engine) {
	env := v1.NewEnv()
	repository := v1.NewRepository()
	chart := v1.NewChart()
	release := v1.NewRelease()

	// 注册中间件
	if global.MyHelmConfig.RunMode == "debug" {
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
	} else {
		router.Use(middleware.AccessLog())
		router.Use(middleware.Recovery())
	}
	router.Use(middleware.ContextTimeout(global.MyHelmConfig.DefaultContextTimeout * time.Second))
	router.Use(middleware.Translations())
	router.Use(middleware.AppInfo())

	// 注册监控接口
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)

	// 注册业务接口
	apiv1 := router.Group("/api/v1")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
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
