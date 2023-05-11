package json

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Upsert : upsert thông tin json, script for mobile
// {"name": "Jasmine"}
func (r Repo) Upsert(data string) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	condition := bson.M{}
	condition["name"] = "Jasmine"

	update := bson.M{
		"$set": bson.M{
			"data":       data,
			"created_at": now,
		},
	}

	_, err = collection.UpdateOne(ctx,
		condition,
		update,
		options.Update().SetUpsert(true))

	return
}
