package clone

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// All : Lấy danh sách clone mà device có
func (r Repo) AllByDevice(req model.AllDeviceReq) (total int, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set filter
	filter := bson.M{}
	filter["alive_status"] = "stored"

	if req.DeviceID != "" {
		filter["device_id"] = req.DeviceID
	}

	//Đếm số lượn clone mà device có
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	total = int(count)

	return
}
