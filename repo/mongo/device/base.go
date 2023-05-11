package device

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
}

var (
	instance *Repo
	once     sync.Once
)

// New ..
func New(ctx context.Context) *Repo {
	once.Do(func() {
		instance = &Repo{
			repo.MongoV2{
				Session:    config.GetConfig().Mongo.Get("core"),
				Collection: "devices",
			},
		}
		collection := instance.Session.GetCollection(instance.Collection)
		{
			indexOpts := options.Index()
			indexOpts.SetBackground(true)

			indexKey := mongo.IndexModel{
				Keys:    bsonx.Doc{{Key: "token", Value: bsonx.Int32(1)}},
				Options: indexOpts,
			}
			_, err := collection.Indexes().CreateOne(ctx, indexKey)

			if err != nil {
				panic(err)
			}
		}
		{
			indexOpts := options.Index()
			indexOpts.SetBackground(true)

			indexKey := mongo.IndexModel{
				Keys:    bsonx.Doc{{Key: "device_id", Value: bsonx.Int32(1)}},
				Options: indexOpts,
			}
			_, err := collection.Indexes().CreateOne(ctx, indexKey)

			if err != nil {
				panic(err)
			}
		}
		{
			indexOpts := options.Index()
			indexOpts.SetBackground(true)

			indexKey := mongo.IndexModel{
				Keys:    bsonx.Doc{{Key: "MacAddress", Value: bsonx.Int32(1)}},
				Options: indexOpts,
			}
			_, err := collection.Indexes().CreateOne(ctx, indexKey)

			if err != nil {
				panic(err)
			}
		}
	})
	// re-init context
	instance.InitContext(ctx)
	return instance
}
