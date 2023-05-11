package service

import (
	"app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetOneRandomActiveType : Lấy 1 service chưa hoàn thành
// . status: Active
// . sort get_random_time : get random
func (r Repo) GetOneRandomActiveType(serviceType string, whitelistTokens []string) (result model.Service, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.FindOneAndUpdate()

	// Sort default theo get_random_time
	option.SetSort(bson.M{"_id": -1, "get_random_time": 1})

	// set condition
	condition := bson.M{}
	condition["type"] = serviceType
	condition["status"] = "Active"
	//! Query temp theo token, sau này remove khi thêm gói
	if len(whitelistTokens) > 0 {
		condition["token"] = bson.M{
			"$in": whitelistTokens,
		}
	}
	fromID, _ := primitive.ObjectIDFromHex("60be39abf0b5761b1e60683d")
	condition["_id"] = bson.M{"$gt": fromID}
	// !=============================================================

	// set update
	update := bson.M{
		"$set": bson.M{
			"get_random_time": time.Now(),
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update, option).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
