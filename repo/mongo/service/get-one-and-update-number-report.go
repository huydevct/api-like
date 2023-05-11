package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetOneAndUpsertNumberReport : Tìm 1 phần tử có
// 1. Tìm theo service_code
// 2. Tăng số number report
func (r Repo) GetOneAndUpsertNumberReport(serviceCode string) (err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	// overwrite data
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["service_code"] = serviceCode

	// set update
	update := bson.M{
		"$set": bson.M{
			"updated_date": now,
			"update_by":    "service_report",
		},
		"$inc": bson.M{
			"number_report": 1,
		},
	}

	myResult := collection.FindOneAndUpdate(
		ctx,
		condition,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	if myResult.Err() != nil {
		err = myResult.Err()
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}

	return
}
