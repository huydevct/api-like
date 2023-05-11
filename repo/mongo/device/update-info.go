package device

import (
	"app/constants"
	"app/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateInfo : Cập nhật info của device
func (r Repo) UpdateInfo(req model.CloneInfo) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["_id"] = req.ID
	condition["status"] = bson.M{"$ne": constants.Delete}

	// set update
	update := bson.M{}
	update["updated_date"] = time.Now()
	if req.Birthday != "" {
		update["birthday"] = req.Birthday
	}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.PhoneNumber != "" {
		update["phone_number"] = req.PhoneNumber
	}
	if req.Sex != "" {
		update["sex"] = req.Sex
	}
	if req.Follow > 0 {
		update["follow"] = req.Follow
	}
	if req.Friend > 0 {
		update["friend"] = req.Friend
	}

	updates := bson.M{
		"$set": update,
	}

	updateResult, err := collection.UpdateOne(ctx,
		condition,
		updates)
	if err != nil {
		return
	}
	if updateResult.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}
