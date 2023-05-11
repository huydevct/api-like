package clone

import (
	"app/constants"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateManyLiveWithSessionCtx : chuyển trạng thái clone thuộc device sang live
// 1. condition: AndoirID
// 2. cập nhật: status: "live"
func (r Repo) UpdateManyLiveWithSessionCtx(sessionCtx context.Context, token string, pcName string) (err error) {

	collection := r.Session.GetCollection(r.Collection)

	// cập nhật data
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["pv_name"] = pcName
	condition["token"] = token
	condition["alive_status"] = constants.CloneStored
	// set update
	update := bson.M{
		"$set": bson.M{
			"alive_status": constants.CloneLive,
			"pc_name":      "",
			"updated_date": now,
		},
	}

	_, err = collection.UpdateMany(sessionCtx,
		condition,
		update)
	if err != nil {
		return
	}
	return
}
