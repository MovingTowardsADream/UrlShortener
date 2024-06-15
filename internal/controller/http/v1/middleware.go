package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		c.Writer.Header().Set("X-Request-ID", requestId)

		c.Set(contextKeyRequestId, requestId)
		c.Next()
	}
}

func Logging(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next() // Обработка запроса

		latency := time.Since(t)
		requestId, _ := c.Get(contextKeyRequestId) // Получение RequestId
		l.Info(fmt.Sprintf("status: %s, url: %s, request_id: %s, status: %d, latency: %s", c.Request.Method, c.Request.URL, requestId, c.Writer.Status(), latency))
	}
}
