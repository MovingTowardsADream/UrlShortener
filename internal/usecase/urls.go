package usecase

import (
	"UrlShortener/internal/entity"
	"UrlShortener/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	_defaultTimeout = 5 * time.Second
)

type UrlsUseCase struct {
	repo    repository.UrlsData
	timeout time.Duration
}

func NewUrlsUseCase(repo repository.UrlsData, opts ...Option) *UrlsUseCase {
	uc := &UrlsUseCase{
		repo:    repo,
		timeout: _defaultTimeout,
	}

	for _, opt := range opts {
		opt(uc)
	}

	return uc
}

func (uc *UrlsUseCase) SaveUrl(ctx context.Context, url entity.Url) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	if url.ShortUrl == "" {
		url.ShortUrl = generateShortLink(url.Url)
	}
	return uc.repo.SaveUrl(ctxWithTimeout, url)
}

func (uc *UrlsUseCase) GetUrl(ctx context.Context, alias string) (entity.Url, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	return uc.repo.GetUrl(ctxWithTimeout, alias)
}

func (uc *UrlsUseCase) DeleteUrl(ctx context.Context, alias string) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	return uc.repo.DeleteUrl(ctxWithTimeout, alias)
}

func generateShortLink(longURL string) string {
	hasher := sha256.New()
	hasher.Write([]byte(longURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}
