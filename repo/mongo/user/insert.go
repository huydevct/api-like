package user

import (
	"context"
	"time"

	"app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert ..
func (r Repo) Insert(data *model.UserInfo) (result *model.UserInfo, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.ID = primitive.NewObjectID()
	data.UpdatedDate = &now
	data.LastChangePasswordDate = &now
	//
	data.CreatedIP = data.UpdatedIP
	data.CreatedUser = data.UpdatedUser
	data.CreatedEmployee = data.UpdatedEmployee
	data.CreatedDate = &now
	data.CreatedSource = data.UpdatedSource
	//
	data.CreatedAt = int(now.UnixNano() / int64(time.Millisecond))

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	result = data
	return
}
