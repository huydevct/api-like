package user

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetActiveByPhone ..
func (r Repo) GetActiveByPhone(phone string) (result model.UserInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{
		"username": phone,
		"status": bson.M{
			"$ne": "Delete",
		}}
	err = collection.FindOne(ctx, condition).Decode(&result)

	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
