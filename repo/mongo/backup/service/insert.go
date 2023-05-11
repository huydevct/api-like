package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"app/model"
)

// Insert ..
func (r Repo) Insert(data model.ActiveService) (result model.ActiveService, err error) {

	ctx := r.GetContext()
	//
	collection := r.Session.GetCollection(r.Collection)
	data.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	result = data
	return
}
