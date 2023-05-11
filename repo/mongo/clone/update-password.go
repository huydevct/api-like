package clone

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdatePassword : cập nhật password, cờ mã hoá
// 1. condition: id, token
// 2. cập nhật: passworod, mz, cz
func (r Repo) UpdatePassword(id primitive.ObjectID, token, password string, mobileEncode, CGIEncode bool) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// cập nhật data
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["_id"] = id
	condition["token"] = token

	// set update
	update := bson.M{
		"$set": bson.M{
			"password":     password,
			"mzz.p":        mobileEncode,
			"czz.p":        CGIEncode,
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
