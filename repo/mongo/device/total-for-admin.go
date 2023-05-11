package device

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// TotalForAdmin : Lấy tổng
func (r Repo) TotalForAdmin(req model.AllDeviceReq) (total int, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set filter
	filter := bson.M{}
	if req.Token != "" {
		filter["token"] = req.Token
	}
	if len(req.Status) > 0 {
		filter["status"] = bson.M{"$in": req.Status}
	}
	if req.DeviceID != "" {
		filter["device_id"] = req.DeviceID
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	total = int(count)

	return
}
