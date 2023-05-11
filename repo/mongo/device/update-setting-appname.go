package device

import (
	"app/constants"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateSettingAppNameReq ..
type UpdateSettingAppNameReq struct {
	DeviceID string
	AppName  constants.AppName
}

// UpdateSettingAppName : Cập nhật appname cho device
func (r Repo) UpdateSettingAppName(req UpdateSettingAppNameReq) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["device_id"] = req.DeviceID
	condition["status"] = bson.M{"$ne": constants.Delete}

	// set update
	update := bson.M{}
	update["updated_date"] = time.Now()
	if req.AppName != "" {
		update["appname"] = req.AppName
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
