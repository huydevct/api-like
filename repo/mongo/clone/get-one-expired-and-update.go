package clone

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneExpiredAndUpdate :
// 1. Tìm 1 phần tử có trạng thái "Getting" và updated_date + 24h < now
func (r Repo) GetOneExpiredAndUpdate() (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()
	expiredDate := now.Add(-24 * time.Hour)

	condition := bson.M{}
	condition["alive_status"] = constants.CloneGetting
	condition["updated_date"] = bson.M{"$lt": expiredDate}

	update := bson.M{"$set": bson.M{
		"alive_status": constants.CloneLive,
		"android_id":   "",
		"updated_date": now,
	}}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
