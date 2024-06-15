package repository

import (
	"UrlShortener/internal/entity"
	"UrlShortener/internal/repository/redisdb"
	"UrlShortener/pkg/redis"
	"context"
)

type UrlsData interface {
	SaveUrl(ctx context.Context, url entity.Url) error
	GetUrl(ctx context.Context, alias string) (entity.Url, error)
	DeleteUrl(ctx context.Context, alias string) error
}

type Repository struct {
	UrlsData
}

func NewRepository(rs *redisdb.Redis) *Repository {
	return &Repository{
		UrlsData: redis_repository.NewUrlsRepo(rs),
	}
}
