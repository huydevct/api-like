package servicelog

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// RemoveMany : RemoveMany by serviceCode
func (r Repo) RemoveMany(serviceCode string) (err error) {

	collection := r.Session.GetCollection(r.Collection)
	ctx := context.Background()

	// set condition
	condition := bson.M{"service_code": serviceCode}

	// set remove
	_, err = collection.DeleteMany(ctx, condition)
	if err != nil {
		return
	}

	return
}
