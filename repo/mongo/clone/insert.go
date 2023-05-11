package clone

import (
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert ..
func (r Repo) Insert(data model.CloneInfo) (result model.CloneInfo, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.ID = primitive.NewObjectID()
	data.CreatedDate = &now
	data.UpdatedDate = &now
	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	result = data
	return
}
