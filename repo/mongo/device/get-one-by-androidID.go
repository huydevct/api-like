package device

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByDeviceID ..
func (r Repo) GetOneByDeviceID(DeviceID string) (result model.DeviceInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{}
	condition["device_id"] = DeviceID
	condition["status"] = bson.M{"$ne": constants.Delete}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneApprovedByDeviceID ..
func (r Repo) GetOneApprovedByDeviceID(DeviceID string) (result model.DeviceInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{}
	condition["DeviceID"] = DeviceID
	condition["status"] = constants.Approved

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
