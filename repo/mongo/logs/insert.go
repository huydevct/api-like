package logs

import (
	"context"
)

// Insert ..
func (r Repo) Insert(data interface{}) (err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	return
}
