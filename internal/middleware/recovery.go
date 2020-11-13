package middleware

import (
	"github.com/gin-gonic/gin"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
	"helm-wrapper/pkg/errcode"
)

func Recovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.WithCallersFrames().Errorf(context, "panic recover err: %v", err)
				app.NewResponse(context).ToErrorResponse(errcode.ServerError)
				context.Abort()
			}
		}()
		context.Next()
	}
}
