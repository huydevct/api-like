package device

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOneAndUpdateLastActionAt : Cập nhật last action at cho DeviceID
func (r Repo) FindOneAndUpdateLastActionAt(DeviceID string) (result model.DeviceInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["device_id"] = DeviceID
	condition["status"] = constants.Approved

	// set update
	update := bson.M{
		"$set": bson.M{
			"last_action_at": now,
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
