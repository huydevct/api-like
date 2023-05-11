package clone

import (
	"app/constants"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateLive : chuyển trạng thái -> live
// 1. condition: cloneID
// 2. cập nhật: status: "live"
func (r Repo) UpdateLive(id primitive.ObjectID) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["_id"] = id

	// set update
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"alive_status": constants.CloneLive,
			"updated_date": now,
		},
	}

	updateResult, err := collection.UpdateOne(ctx,
		condition,
		update)
	if err != nil {
		return
	}
	if updateResult.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}
