package user

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByToken ..
func (r Repo) GetOneByToken(token string) (result model.UserInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{
		"token": token,
		"status": bson.M{
			"$in": []constants.CommonStatus{
				constants.Active,
				constants.Pause,
			},
		}}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneActiveByToken ..
func (r Repo) GetOneActiveByToken(token string) (result model.UserInfo, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"token": token, "status": constants.Active}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
