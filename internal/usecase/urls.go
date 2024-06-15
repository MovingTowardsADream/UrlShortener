package usecase

import (
	"UrlShortener/internal/entity"
	"UrlShortener/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type UrlsUseCase struct {
	repo repository.UrlsData
}

func NewUrlsUseCase(repo repository.UrlsData) *UrlsUseCase {
	return &UrlsUseCase{repo: repo}
}

func (uc *UrlsUseCase) SaveUrl(ctx context.Context, url entity.Url) error {
	if url.ShortUrl == "" {
		url.ShortUrl = generateShortLink(url.Url)
	}
	return uc.repo.SaveUrl(ctx, url)
}

func (uc *UrlsUseCase) GetUrl(ctx context.Context, alias string) (entity.Url, error) {
	return uc.repo.GetUrl(ctx, alias)
}

func (uc *UrlsUseCase) DeleteUrl(ctx context.Context, alias string) error {
	return uc.repo.DeleteUrl(ctx, alias)
}

func generateShortLink(longURL string) string {
	hasher := sha256.New()
	hasher.Write([]byte(longURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}
