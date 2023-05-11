package clone

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOneLiveAndUpdate : Tìm 1 clone có status: live và chuyển thành status: getting
// 1. token
// 2. appname
// 3. update status to "getting"
func (r Repo) FindOneLiveAndUpdate(token string, appName constants.AppName) (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["token"] = token
	condition["appname"] = appName
	// condition["is_autofarmer"] = false
	condition["alive_status"] = constants.CloneLive

	// set update
	update := bson.M{
		"$set": bson.M{
			"alive_status": constants.CloneGetting,
			"updated_date": now,
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// FindOneLiveAndUpdate : Tìm 1 clone có status: live và chuyển thành status: getting
// 1. token
// 2. appname
// 3. update status to "getting"
type FindCloneLiveByDeviceReq struct {
	Token   string
	AppName string
	PCName  string
}

func (r Repo) FindOneLiveAndUpdateByRequest(cloneRequest FindCloneLiveByDeviceReq) (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()
	option := options.FindOneAndUpdate()
	option.SetSort(bson.M{"get_random_time": 1})
	condition := bson.M{}

	condition["token"] = cloneRequest.Token
	condition["alive_status"] = constants.CloneLive
	condition["appname"] = cloneRequest.AppName
	// set update
	update := bson.M{
		"$set": bson.M{
			"pc_name":         cloneRequest.PCName,
			"updated_date":    now,
			"get_random_time": now,
		},
	}

	err = collection.FindOneAndUpdate(ctx, condition, update, option).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
