package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"phos.cc/yoo/internal/pkg/known"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否包含 X-Request-ID 字段，如果没有则生成一个
		requestID := c.Request.Header.Get(known.XRequestIDKey)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将 X-Request-ID 保存到上下文中
		c.Set(known.XRequestIDKey, requestID)

		c.Writer.Header().Set(known.XRequestIDKey, requestID)

		c.Next()
	}
}
