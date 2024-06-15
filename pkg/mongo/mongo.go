package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	defaultMaxPoolSize = 100
	defaultMinIdleConn = 10
	defaultConnTimeout = 5 * time.Second
	defaultDiscTimeout = 5 * time.Second
)

type Mongo struct {
	maxPoolSize int
	minIdleConn int
	connTimeout time.Duration

	Client *mongo.Client
}

func NewMongoClient(url string, opts ...Option) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()

	mg := &Mongo{
		maxPoolSize: defaultMaxPoolSize,
		minIdleConn: defaultMinIdleConn,
		connTimeout: defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(mg)
	}

	clientOptions := options.Client().
		ApplyURI(url).
		SetMaxPoolSize(uint64(mg.maxPoolSize)).
		SetMinPoolSize(uint64(mg.minIdleConn))

	var err error
	mg.Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = mg.Client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return mg, nil
}

func (m *Mongo) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDiscTimeout)
	defer cancel()

	return m.Client.Disconnect(ctx)
}
