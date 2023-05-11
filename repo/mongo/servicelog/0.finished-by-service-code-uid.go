package servicelog

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FinishedByServiceCodeUID :
func (r Repo) FinishedByServiceCodeUID(serviceCode, cloneUID, serviceType, DeviceID string) (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	condition := bson.M{}
	condition["service_code"] = serviceCode
	condition["uid"] = cloneUID
	condition["type"] = serviceType

	update := bson.M{"$set": bson.M{
		"status":       constants.ServicelogFinished,
		"device_id":    DeviceID,
		"updated_date": now,
	}}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
