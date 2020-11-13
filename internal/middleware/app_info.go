package middleware

import "github.com/gin-gonic/gin"

func AppInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("app_name", "helm-wrapper")
		context.Set("app_version", "v1.1")
		context.Next()
	}
}
