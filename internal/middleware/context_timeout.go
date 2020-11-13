package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func ContextTimeout(t time.Duration) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// 构建超时上下文
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), t)
		defer cancel()

		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
