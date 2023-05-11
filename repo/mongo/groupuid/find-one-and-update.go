package groupuid

import (
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOneAndUpdate ..
func (r Repo) FindOneAndUpdate(group_id *primitive.ObjectID) (result model.GroupUIDInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()
	option := options.FindOneAndUpdate()
	// set condition
	condition := bson.M{}
	condition["group_id"] = group_id
	condition["status"] = "Active"
	option.SetSort(bson.M{"get_time_random": 1})

	// set update
	update := bson.M{
		"$set": bson.M{
			"get_time_random": now,
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update, option).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}
	return
}
