package service

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetOneRandomService : Tìm n phần tử có
// 1. type
// 2. number_rest > 0
func (r Repo) GetOneRandomService(serviceType string) (result model.ActiveService, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	// Sort default theo get_random_time
	// set option go mod download

	option := options.FindOneAndUpdate()
	option.SetSort(bson.M{"get_random_time": 1})

	// set condition
	condition := bson.M{}
	condition["status"] = constants.ServiceActive
	condition["type"] = serviceType

	// Close the cursor once finished
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	// set update
	update := bson.M{
		"$inc": bson.M{"number_rest": -1},
		"$set": bson.M{
			"get_random_time": time.Now(),
			"updated_date":    time.Now(),
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update, option).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
