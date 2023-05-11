package device

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateDebugModeReq ..
type UpdateDebugModeReq struct {
	DeviceID  primitive.ObjectID
	DebugMode *int
}

// UpdateDebugMode : Cập nhật info của device
func (r Repo) UpdateDebugMode(req UpdateDebugModeReq) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["_id"] = req.DeviceID

	// set update
	update := bson.M{}
	update["updated_date"] = time.Now()

	if req.DebugMode != nil {
		update["debug_mode"] = *req.DebugMode
	}

	updates := bson.M{
		"$set": update,
	}

	updateResult, err := collection.UpdateOne(ctx,
		condition,
		updates)
	if err != nil {
		return
	}
	if updateResult.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}
