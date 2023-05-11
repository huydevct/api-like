package token

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Detail ..
func (r Repo) Detail(condition bson.M) (result model.Token, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
