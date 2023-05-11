package clone

import (
	"app/constants"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateChecking : chuyển trạng thái -> checking
// 1. condition: cloneID, token
// 2. cập nhật: status: "checking"
func (r Repo) UpdateCheckingToLive(id primitive.ObjectID, status constants.AliveStatus) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["_id"] = id

	// set update
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"alive_status": status,
			"android_id":   "",
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
