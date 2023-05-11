package clone

import (
	"app/constants"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateCheckpoint : chuyển trạng thái "live", "getting" -> checkpoint
// 1. condition: cloneID, token, status: ["live", "getting"]
// 2. cập nhật: status: "stored", DeviceID
func (r Repo) UpdateCheckpoint(id primitive.ObjectID, token, pcName string) (result *mongo.UpdateResult, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// cập nhật data
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["_id"] = id
	condition["token"] = token
	condition["alive_status"] = bson.M{
		"$in": []constants.AliveStatus{
			constants.CloneLive,
			constants.CloneGetting,
			constants.CloneStored,
		},
	}

	// set update
	update := bson.M{
		"$set": bson.M{
			"alive_status": constants.CloneCheckpoint,
			"pc_name":      pcName,
			"updated_date": now,
		},
		"$addToSet": bson.M{
			"checkpoint_device_ids": pcName,
		},
	}

	result, err = collection.UpdateOne(ctx,
		condition,
		update)
	if err != nil {
		return
	}
	if result.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}
