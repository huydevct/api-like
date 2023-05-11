package user

import (
	"context"
	"fmt"
	"time"

	"app/model"

	"go.mongodb.org/mongo-driver/bson"
)

// Update ..
func (r Repo) Update(data model.UserInfo) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.UpdatedDate = &now
	condition := bson.M{"_id": data.ID}

	updateResult, err := collection.UpdateOne(ctx,
		condition,
		bson.M{"$set": data})
	if err != nil {
		return
	}
	if updateResult.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}

// UpdateWithSessionCtx : cập nhật user với transaction
func (r Repo) UpdateWithSessionCtx(sessionCtx context.Context, data model.UserInfo) (err error) {

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.UpdatedDate = &now
	condition := bson.M{"_id": data.ID}

	updateResult, err := collection.UpdateOne(sessionCtx,
		condition,
		bson.M{"$set": data})
	if err != nil {
		return
	}
	if updateResult.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}
