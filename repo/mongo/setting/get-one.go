package setting

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOne ..
func (r Repo) GetOne() (result model.Setting, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	err = collection.FindOne(ctx, bson.M{}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
