package token

import (
	"context"
	"time"

	"app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create func
func (r Repo) Create(data model.Token) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	// Update data
	now := time.Now()
	expired := now.Add(time.Second * time.Duration(data.ExpiredAfterSecond))
	//
	data.ID = primitive.NewObjectID()
	data.CreatedDate = &now
	data.ExpiredDate = &expired

	_, err = collection.InsertOne(ctx, data)
	return
}
