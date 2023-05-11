package clone

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByEmail ..
func (r Repo) GetOneByEmail(email string) (result model.CloneInfo, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"email": email}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneActiveByEmail ..
func (r Repo) GetOneActiveByEmail(email string) (result model.CloneInfo, err error) {

	ctx := context.Background()
	//
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"email": email, "alive_status": bson.M{"$ne": constants.CloneDelete}}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
