package service

import (
	"context"
	"sync"

	"app/common/config"
	"app/repo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Repo type
type Repo struct {
	repo.MongoV2
	ctx context.Context
}

// GetContext get global context when init
func (r *Repo) GetContext() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

var (
	instance Repo
	once     sync.Once
	mu       sync.RWMutex
)

// New ..
func New(ctx context.Context) (result Repo) {
	mu.Lock()

	once.Do(func() {
		instance = Repo{
			repo.MongoV2{
				Session:    config.GetConfig().Mongo.Get("autolike"),
				Collection: "services",
			},
			ctx,
		}
		collection := instance.Session.GetCollection(instance.Collection)
		{
			indexOpts := options.Index()
			indexOpts.SetBackground(true)

			indexKey := mongo.IndexModel{
				Keys:    bsonx.Doc{{Key: "service_code", Value: bsonx.Int32(1)}},
				Options: indexOpts,
			}
			_, err := collection.Indexes().CreateOne(ctx, indexKey)

			if err != nil {
				panic(err)
			}
		}
	})
	// Tạo 1 bảng copy biến instance
	result = instance
	result.ctx = ctx

	mu.Unlock()
	return

}
