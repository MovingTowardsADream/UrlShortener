package mongo_repository

import (
	"UrlShortener/internal/entity"
	mongodb "UrlShortener/pkg/mongo"
	"context"
)

type UrlsRepo struct {
	db *mongodb.Mongo
}

func NewUrlsRepo(mg *mongodb.Mongo) *UrlsRepo {
	return &UrlsRepo{mg}
}

func (ur *UrlsRepo) SaveUrl(ctx context.Context, url entity.Url) error {
	return nil
}

func (ur *UrlsRepo) GetUrl(ctx context.Context, alias string) (entity.Url, error) {
	var url entity.Url
	return url, nil
}

func (ur *UrlsRepo) DeleteUrl(ctx context.Context, alias string) error {
	return nil
}
