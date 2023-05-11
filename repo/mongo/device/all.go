package device

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách device
func (r Repo) All(req model.AllDeviceReq) (results []model.DeviceInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"_id": -1})

	if req.Limit != -1 {
		if req.Limit > 0 && req.Limit <= 100 {
			option.SetLimit(int64(req.Limit))
		} else {
			option.SetLimit(100)
		}
	}

	if req.Page >= 0 {
		if req.Page == 0 {
			req.Page = 1
		}
		skip := (req.Page - 1) * req.Limit

		option.SetSkip(int64(skip))
	}

	// set filter
	filter := bson.M{}
	if !req.Offset.IsZero() {
		filter["_id"] = bson.M{"$lt": req.Offset}
	}
	if req.Token != "" {
		filter["token"] = req.Token
	}
	if len(req.Status) > 0 {
		filter["status"] = bson.M{"$in": req.Status}
	}
	if req.DeviceID != "" {
		filter["device_id"] = req.DeviceID
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.DeviceInfo{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
