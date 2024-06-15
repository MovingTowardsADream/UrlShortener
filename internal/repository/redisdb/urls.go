package redis_repository

import (
	"UrlShortener/internal/entity"
	"UrlShortener/pkg/redis"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type UrlsRepo struct {
	db *redisdb.Redis
}

func NewUrlsRepo(rs *redisdb.Redis) *UrlsRepo {
	return &UrlsRepo{rs}
}

func (ur *UrlsRepo) SaveUrl(ctx context.Context, url entity.Url) error {
	err := ur.db.Client.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Get(ctx, url.ShortUrl).Result()

		if err != redis.Nil {
			return fmt.Errorf("UrlsRepo.SaveUrl - tx.Get: %v", err)
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, url.ShortUrl, url.Url, url.Expiry)
			return nil
		})
		return err
	})

	if err != nil {
		return fmt.Errorf("UrlsRepo.SaveUrl - Watch: %v", err)
	}
	return nil
}

func (ur *UrlsRepo) GetUrl(ctx context.Context, alias string) (entity.Url, error) {
	var value string
	var ttl time.Duration

	err := ur.db.Client.Watch(ctx, func(tx *redis.Tx) error {
		var err error
		value, err = ur.db.Client.Get(ctx, alias).Result()
		if err != nil {
			return fmt.Errorf("UrlsRepo.GetUrl - Get: %v", err)
		}

		ttl, err = ur.db.Client.TTL(ctx, alias).Result()
		if err != nil {
			return fmt.Errorf("UrlsRepo.GetUrl - Get: %v", err)
		}
		return nil
	})

	var url entity.Url

	if err != nil {
		return url, fmt.Errorf("UrlsRepo.GetUrl - Watch: %v", err)
	}

	url = entity.Url{
		Url:      value,
		ShortUrl: alias,
		Expiry:   ttl,
	}

	return url, nil
}

func (ur *UrlsRepo) DeleteUrl(ctx context.Context, alias string) error {
	err := ur.db.Client.Del(ctx, alias).Err()
	if err != nil {
		return fmt.Errorf("UrlsRepo.DeleteUrl - Del: %v", err)
	}
	return nil
}
