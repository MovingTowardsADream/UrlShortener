package v1

import (
	"UrlShortener/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	defaultLogsFile = "./logs/requests.log"
)

const (
	contextKeyRequestId = "requestId"
)

func NewRouter(handler *gin.Engine, l *slog.Logger, u *usecase.UseCase) {
	//handler.Use(RequestId())

	handler.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf(`{"time":"%s", "method":"%s","uri":"%s", "status":%d, "latency":%s,"error":"%s"}`+"\n",
				param.TimeStamp.Format(time.RFC3339Nano),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				param.ErrorMessage,
			)
		},
		Output: io.MultiWriter(os.Stdout, setLogsFile()),
	}))

	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/api/v1")
	{
		newUrlsRoutes(h, u.UrlsData, l)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile(defaultLogsFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
