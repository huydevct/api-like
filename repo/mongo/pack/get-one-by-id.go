package pack

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByID ..
func (r Repo) GetOneByID(id primitive.ObjectID) (result model.Package, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{
		"_id": id,
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

// GetOneActiveByID ..
func (r Repo) GetOneActiveByID(id primitive.ObjectID) (result model.Package, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"_id": id, "status": constants.Active}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
