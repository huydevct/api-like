package servicelog

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneCreatedAndUpdate :
// 1. Tìm 1 phần tử có trạng thái "Created"
// 2. Update trạng thái phần tử này thành "Updating" để tránh bị update trùng
func (r Repo) GetOneCreatedAndUpdate(serviceCode string) (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	condition := bson.M{}
	condition["status"] = constants.ServicelogCreated
	condition["updated_date"] = bson.M{"$lt": now}
	if serviceCode != "" {
		condition["service_code"] = serviceCode
	}

	update := bson.M{"$set": bson.M{
		"status":       constants.ServicelogUpdating,
		"updated_date": now,
	}}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
