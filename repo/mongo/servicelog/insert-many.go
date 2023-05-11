package servicelog

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"app/model"
)

// InsertMany ..
func (r Repo) InsertMany(data []model.ServiceLog) (err error) {

	collection := r.Session.GetCollection(r.Collection)
	ctx := context.Background()

	// Cạp nhật Data
	insertData := make([]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		data[i].ID = primitive.NewObjectID()
		//
		insertData[i] = data[i]
	}

	_, err = collection.InsertMany(ctx,
		insertData,
		options.InsertMany().SetOrdered(false))

	return
}
