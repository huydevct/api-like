package friend

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByID ..
func (r Repo) GetOneByID(id primitive.ObjectID) (result model.FriendInfo, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"_id": id}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneActiveByID ..
func (r Repo) GetOneActiveByID(id primitive.ObjectID) (result model.FriendInfo, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{"_id": id}
	condition["_id"] = id
	condition["status"] = constants.Active

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
