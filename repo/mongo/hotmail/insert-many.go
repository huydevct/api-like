package hotmail

import (
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertMany ..
func (r Repo) InsertMany(data []model.HotMail) (hotmail []model.HotMail, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	now := time.Now()
	// Cạp nhật Data
	insertData := make([]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		data[i].ID = primitive.NewObjectID()
		//
		data[i].CreatedDate = &now
		insertData[i] = data[i]
		hotmail = append(hotmail, data[i])
	}

	_, err = collection.InsertMany(ctx, insertData, options.InsertMany().SetOrdered(false))
	if err != nil {
		return
	}
	return
}
