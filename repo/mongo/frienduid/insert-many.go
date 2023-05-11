package frienduid

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"app/model"
)

// InsertMany ..
func (r Repo) InsertMany(data []model.FriendUIDInfo) (err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)

	// Cạp nhật Data
	insertData := make([]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		data[i].ID = primitive.NewObjectID()
		insertData[i] = data[i]
	}

	_, err = collection.InsertMany(ctx, insertData, options.InsertMany().SetOrdered(false))
	if err != nil {
		return
	}

	return
}

// InsertManyWithSessionCtx : Thực hiện insertMany với transaction
func (r Repo) InsertManyWithSessionCtx(sessionCtx context.Context, data []model.FriendUIDInfo) (err error) {

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// Cạp nhật Data
	insertData := make([]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		data[i].ID = primitive.NewObjectID()
		data[i].CreatedDate = &now

		insertData[i] = data[i]
	}

	_, err = collection.InsertMany(sessionCtx, insertData)
	if err != nil {
		return
	}

	return
}
