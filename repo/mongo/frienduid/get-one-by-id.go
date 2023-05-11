package frienduid

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByID ..
func (r Repo) GetOneByID(id primitive.ObjectID) (result model.FriendUIDInfo, err error) {

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
