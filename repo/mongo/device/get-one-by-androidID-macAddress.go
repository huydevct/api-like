package device

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByDeviceIDMacAddress ..
func (r Repo) GetOneByDeviceIDMacAddress(DeviceID, macAddress string) (result model.DeviceInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{}
	condition["device_id"] = DeviceID
	condition["MacAddress"] = macAddress
	condition["status"] = bson.M{"$ne": constants.Delete}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
