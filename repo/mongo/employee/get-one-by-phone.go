package employee

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByUsername ..
func (r Repo) GetOneByUsername(username string) (result model.UserInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{
		"username": username,
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
