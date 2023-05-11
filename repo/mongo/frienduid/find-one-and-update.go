package frienduid

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
func (r Repo) FindOneAndUpdate(friend_id *primitive.ObjectID) (result model.FriendUIDInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["friend_id"] = friend_id
	condition["status"] = "Active"
	option := options.FindOneAndUpdate()
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
