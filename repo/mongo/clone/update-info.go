package clone

import (
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateInfo ..
func (r Repo) UpdateInfo(req model.CloneInfo) (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["_id"] = req.ID

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

	myResult := collection.FindOneAndUpdate(ctx,
		condition,
		updates,
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	if myResult.Err() != nil {
		err = myResult.Err()
		return
	}

	err = myResult.Decode(&result)
	if err != nil {
		return
	}

	return
}
