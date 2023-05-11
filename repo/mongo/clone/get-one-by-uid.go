package clone

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByUID ..
func (r Repo) GetOneByUID(uid string) (result model.CloneInfo, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"uid": uid}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneActiveByUID ..
func (r Repo) GetOneActiveByUID(uid string) (result model.CloneInfo, err error) {

	ctx := context.Background()
	//
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"uid": uid, "alive_status": bson.M{"$ne": constants.CloneDelete}}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
