package redisdb

import "time"

type Option func(*Redis)

func MaxConnTimeout(time time.Duration) Option {
	return func(r *Redis) {
		r.connTimeout = time
	}
}

func MaxPoolSize(maxPool int) Option {
	return func(r *Redis) {
		r.maxPoolSize = maxPool
	}
}

func MinIdleConn(minIdle int) Option {
	return func(r *Redis) {
		r.minIdleConn = minIdle
	}
}
