package gift

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"app/model"
)

// InsertMany ..
func (r Repo) InsertMany(data []model.GiftCode) (err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// Cạp nhật Data
	insertData := make([]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		data[i].ID = primitive.NewObjectID()
		data[i].UpdatedDate = &now
		//
		data[i].CreatedIP = data[i].UpdatedIP
		data[i].CreatedUser = data[i].UpdatedUser
		data[i].CreatedEmployee = data[i].UpdatedEmployee
		data[i].CreatedDate = &now
		data[i].CreatedSource = data[i].UpdatedSource
		//
		data[i].CreatedAt = int(now.UnixNano() / int64(time.Millisecond))

		insertData[i] = data[i]
	}

	_, err = collection.InsertMany(ctx, insertData, options.InsertMany().SetOrdered(false))
	if err != nil {
		return
	}

	return
}

// InsertManyWithSessionCtx : Thực hiện insertMany với transaction
func (r Repo) InsertManyWithSessionCtx(sessionCtx context.Context, data []model.GiftCode) (err error) {

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// Cạp nhật Data
	insertData := make([]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		data[i].ID = primitive.NewObjectID()
		data[i].UpdatedDate = &now
		//
		data[i].CreatedIP = data[i].UpdatedIP
		data[i].CreatedUser = data[i].UpdatedUser
		data[i].CreatedEmployee = data[i].UpdatedEmployee
		data[i].CreatedDate = &now
		data[i].CreatedSource = data[i].UpdatedSource
		//
		data[i].CreatedAt = int(now.UnixNano() / int64(time.Millisecond))

		insertData[i] = data[i]
	}

	_, err = collection.InsertMany(sessionCtx, insertData)
	if err != nil {
		return
	}

	return
}
