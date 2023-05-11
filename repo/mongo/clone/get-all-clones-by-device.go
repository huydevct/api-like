package clone

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AllClonesByDeviceReq struct {
	DeviceID string
	Token    string
}

// GetAllCloneByDevice : Lấy danh sách clone theo device
func (r Repo) GetAllCloneByDevice(req AllClonesByDeviceReq) (results []model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"_id": -1})

	// set filter
	filter := bson.M{}
	if req.Token != "" {
		filter["token"] = req.Token
	}
	if req.DeviceID != "" {
		filter["device_id"] = req.DeviceID
	}
	filter["alive_status"] = constants.CloneStored

	cur, err := collection.Find(
		ctx,
		filter,
		option,
	)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		var temp model.CloneInfo
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
