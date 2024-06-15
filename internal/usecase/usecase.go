package usecase

import (
	"UrlShortener/internal/entity"
	"UrlShortener/internal/repository"
	"context"
)

type UrlsData interface {
	SaveUrl(ctx context.Context, url entity.Url) error
	GetUrl(ctx context.Context, alias string) (entity.Url, error)
	DeleteUrl(ctx context.Context, alias string) error
}

type UseCase struct {
	UrlsData
}

func NewUseCase(repos *repository.Repository) *UseCase {
	return &UseCase{
		UrlsData: NewUrlsUseCase(repos.UrlsData),
	}
}
