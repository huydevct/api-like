package setting

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"app/model"
)

// Insert ..
func (r Repo) Insert(data model.Setting) (result model.Setting, err error) {

	ctx := context.Background()

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
