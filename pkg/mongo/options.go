package mongodb

import "time"

type Option func(*Mongo)

func MaxConnTimeout(time time.Duration) Option {
	return func(r *Mongo) {
		r.connTimeout = time
	}
}

func MaxPoolSize(maxPool int) Option {
	return func(r *Mongo) {
		r.maxPoolSize = maxPool
	}
}

func MinIdleConn(minIdle int) Option {
	return func(r *Mongo) {
		r.minIdleConn = minIdle
	}
}
