package redisdb

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	defaultMaxPoolSize = 100
	defaultMinIdleConn = 10
	defaultConnTimeout = 5 * time.Second
)

type Redis struct {
	maxPoolSize int
	minIdleConn int
	connTimeout time.Duration

	Client *redis.Client
}

func NewRedisClient(addr string, password string, db int, opts ...Option) (*Redis, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()

	rs := &Redis{
		maxPoolSize: defaultMaxPoolSize,
		minIdleConn: defaultMinIdleConn,
		connTimeout: defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(rs)
	}

	rs.Client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     rs.maxPoolSize,
		MinIdleConns: rs.minIdleConn,
	})

	_, err := rs.Client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection error: %w", err)
	}

	return rs, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
