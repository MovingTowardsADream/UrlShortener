package app

import (
	"UrlShortener/configs"
	v1 "UrlShortener/internal/controller/http/v1"
	"UrlShortener/internal/repository"
	"UrlShortener/internal/usecase"
	"UrlShortener/pkg/httpserver"
	redisdb "UrlShortener/pkg/redis"
	"UrlShortener/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log/slog"
)

type App struct {
	HTTPServer *httpserver.Server
	DB         *redisdb.Redis
}

func New(log *slog.Logger, cfg *configs.Config) *App {

	// NoSQL db
	rs, err := redisdb.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password, 0)

	if err != nil {
		panic("app - New - Redis - error redis running: " + err.Error())
	}

	repo := repository.NewRepository(rs)

	useCase := usecase.NewUseCase(repo)

	// Init http server
	handler := gin.New()

	// Custom validator
	binding.Validator = validator.NewCustomValidator()

	v1.NewRouter(handler, log, useCase)
	httpServer := httpserver.New(log, handler, httpserver.Port(cfg.HTTP.Port), httpserver.WriteTimeout(cfg.HTTP.Timeout))

	return &App{
		HTTPServer: httpServer,
		DB:         rs,
	}
}
