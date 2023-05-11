package setting

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneSetting ..
func (r Repo) GetOneSetting() (result model.Setting, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
