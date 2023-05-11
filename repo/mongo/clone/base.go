package clone

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
				Collection: "clones",
			},
		}
		collection := instance.Session.GetCollection(instance.Collection)
		{
			indexOpts := options.Index()
			indexOpts.SetBackground(true)

			indexKey := mongo.IndexModel{
				Keys: bsonx.Doc{
					{Key: "token", Value: bsonx.Int32(1)},
					{Key: "appname", Value: bsonx.Int32(1)},
					{Key: "alive_status", Value: bsonx.Int32(1)},
					{Key: "get_random_time", Value: bsonx.Int32(1)},
				},
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
				Keys:    bsonx.Doc{{Key: "uid", Value: bsonx.Int32(1)}},
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
				Keys:    bsonx.Doc{{Key: "email", Value: bsonx.Int32(1)}},
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
				Keys:    bsonx.Doc{{Key: "android_id", Value: bsonx.Int32(1)}},
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
				Keys: bsonx.Doc{
					{Key: "token", Value: bsonx.Int32(1)},
					{Key: "android_id", Value: bsonx.Int32(1)},
					{Key: "alive_status", Value: bsonx.Int32(1)},
				},
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
				Keys: bsonx.Doc{
					{Key: "alive_status", Value: bsonx.Int32(1)},
					{Key: "updated_date", Value: bsonx.Int32(1)},
				},
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
				Keys: bsonx.Doc{
					{Key: "appname", Value: bsonx.Int32(1)},
					{Key: "alive_status", Value: bsonx.Int32(1)},
					{Key: "get_random_time", Value: bsonx.Int32(1)},
				},
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
				Keys: bsonx.Doc{
					{Key: "appname", Value: bsonx.Int32(1)},
					{Key: "alive_status", Value: bsonx.Int32(1)},
					{Key: "updated_clone_checkpoint", Value: bsonx.Int32(1)},
					{Key: "get_random_time", Value: bsonx.Int32(1)},
				},
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
