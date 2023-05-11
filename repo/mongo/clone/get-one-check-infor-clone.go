package clone

import (
	"app/constants"
	"app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetOneCheckInfoClone : Lấy 1 clone chưa check info
func (r Repo) GetOneCheckInfoClone(aliveStatus constants.AliveStatus) (result model.CloneInfo, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)

	now := time.Now()

	// set option
	option := options.FindOneAndUpdate()

	// Sort default theo get_random_time
	option.SetSort(bson.M{"get_random_time": 1})

	// set condition
	condition := bson.M{}
	condition["alive_status"] = aliveStatus
	condition["appname"] = constants.AppNameFaceBook

	// set update
	update := bson.M{
		"$set": bson.M{
			"get_random_time": now,
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update, option).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
