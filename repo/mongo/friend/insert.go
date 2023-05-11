package friend

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"app/model"
)

// Insert ..
func (r Repo) Insert(data model.FriendInfo) (result model.FriendInfo, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.ID = primitive.NewObjectID()
	data.UpdatedDate = &now
	//
	data.CreatedIP = data.UpdatedIP
	data.CreatedUser = data.UpdatedUser
	data.CreatedEmployee = data.UpdatedEmployee
	data.CreatedDate = &now
	data.CreatedSource = data.UpdatedSource

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	result = data
	return
}

// InsertWithSessionCtx : Thực hiện insert với transaction
func (r Repo) InsertWithSessionCtx(sessionCtx context.Context, data model.FriendInfo) (result model.FriendInfo, err error) {

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.ID = primitive.NewObjectID()
	data.UpdatedDate = &now
	//
	data.CreatedIP = data.UpdatedIP
	data.CreatedUser = data.UpdatedUser
	data.CreatedEmployee = data.UpdatedEmployee
	data.CreatedDate = &now
	data.CreatedSource = data.UpdatedSource

	_, err = collection.InsertOne(sessionCtx, data)
	if err != nil {
		return
	}

	result = data
	return
}
