package user

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByToken ..
func (r Repo) GetTokenNewByOld(token string) (result model.UserInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{
		"old_token": token,
	}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
