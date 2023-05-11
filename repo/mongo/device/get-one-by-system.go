package device

import (
	"context"

	"app/constants"
	"app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByID ..
func (r Repo) GetOneBySystem(Token string, DeviceID string) (result model.DeviceInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{}
	condition["token"] = Token
	condition["device_id"] = DeviceID
	condition["status"] = constants.Active

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
