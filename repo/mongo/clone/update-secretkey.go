package clone

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateSecretkey : cập nhật Secretkey, cờ mã hoá
// 1. condition: id, token
// 2. cập nhật: secretkey, mz, cz
func (r Repo) UpdateSecretkey(id primitive.ObjectID, token, secretkey string, mobileEncode, CGIEncode bool) (err error) {

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
			"secretkey":    secretkey,
			"mzz.s":        mobileEncode,
			"czz.s":        CGIEncode,
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
