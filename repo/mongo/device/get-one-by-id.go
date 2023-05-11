package device

import (
	"context"

	"app/constants"
	"app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByID ..
func (r Repo) GetOneByID(id primitive.ObjectID) (result model.DeviceInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{}
	condition["_id"] = id
	condition["status"] = bson.M{"$ne": constants.Delete}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
