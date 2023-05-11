package service

import (
	"fmt"
	"time"

	"app/model"

	"go.mongodb.org/mongo-driver/bson"
)

// Update ..
func (r Repo) Update(data model.Service) (err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()
	// cập nhật update time
	data.UpdatedDate = &now
	data.UpdateTime = int(now.UnixNano() / int64(time.Millisecond))

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
