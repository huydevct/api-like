package actionprofile

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"app/model"
)

// Insert ..
func (r Repo) Insert(data model.ActionProfile) (result model.ActionProfile, err error) {

	ctx := context.Background()
	//
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.ID = primitive.NewObjectID()

	data.CreatedDate = &now

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	result = data
	return
}
