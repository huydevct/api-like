package clone

import (
	"app/constants"
	"app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetOneCloneCheckpoint : Lấy 1 clone chưa check info
func (r Repo) GetOneCloneCheckpoint() (result model.CloneInfo, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)

	// set time
	now := time.Now()
	then := now.AddDate(0, 0, -3)
	// set option
	option := options.FindOneAndUpdate()

	// Sort default theo get_random_time
	option.SetSort(bson.M{"get_random_time": 1})

	// set condition
	condition := bson.M{}
	condition["alive_status"] = constants.CloneCheckpoint
	condition["appname"] = constants.AppNameFaceBook
	orQuery := []bson.M{}
	orQuery = append(orQuery, bson.M{"updated_clone_checkpoint": bson.M{"$exists": false}}, bson.M{"updated_clone_checkpoint": bson.M{"$lte": then}})
	condition["$or"] = orQuery

	// set update
	update := bson.M{
		"$set": bson.M{
			"get_random_time": time.Now(),
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update, option).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
