package clone

import (
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOneAndUpdateDeviceID : support cho migrate hệ thống, cân nhắc khi dùng
// 1. uid
// 2. token
// 3. update DeviceID
func (r Repo) FindOneAndUpdateDeviceID(uid, token, DeviceID string) (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["uid"] = uid
	condition["token"] = token

	// set update
	update := bson.M{
		"$set": bson.M{
			"device_id":    DeviceID,
			"updated_date": now,
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
