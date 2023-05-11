package service

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"app/model"
)

// Insert ..
func (r Repo) Insert(data model.Service) (result model.Service, err error) {

	ctx := r.GetContext()
	//
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.ID = primitive.NewObjectID()
	data.UpdatedDate = &now
	//
	data.FanpageID = strings.TrimSpace(data.FanpageID)
	data.CreatedIP = data.UpdatedIP
	data.CreatedUser = data.UpdatedUser
	data.CreatedEmployee = data.UpdatedEmployee
	data.CreatedDate = &now
	data.CreatedSource = data.UpdatedSource
	//
	data.CreatedAt = int(now.UnixNano() / int64(time.Millisecond))
	data.UpdateTime = int(now.UnixNano() / int64(time.Millisecond))

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	result = data
	return
}

// InsertWithSessionCtx : update with transaction
func (r Repo) InsertWithSessionCtx(sessionCtx context.Context, data model.Service) (result model.Service, err error) {

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.ID = primitive.NewObjectID()
	data.UpdatedDate = &now
	//
	data.FanpageID = strings.TrimSpace(data.FanpageID)
	data.CreatedIP = data.UpdatedIP
	data.CreatedUser = data.UpdatedUser
	data.CreatedEmployee = data.UpdatedEmployee
	data.CreatedDate = &now
	data.CreatedSource = data.UpdatedSource
	//
	data.CreatedAt = int(now.UnixNano() / int64(time.Millisecond))
	data.UpdateTime = int(now.UnixNano() / int64(time.Millisecond))

	_, err = collection.InsertOne(sessionCtx, data)
	if err != nil {
		return
	}

	result = data
	return
}
