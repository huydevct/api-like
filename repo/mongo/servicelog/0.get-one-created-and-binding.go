package servicelog

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneCreatedAndBindingReq ..
type GetOneCreatedAndBindingReq struct {
	ServiceCode string
	UID         string
	DeviceID    string
	Token       string
	AppName     string
}

// GetOneCreatedAndBinding :
// 1. Tìm 1 phần tử có trạng thái "Created", serviceCode
// 2. update binding với uid
func (r Repo) GetOneCreatedAndBinding(data GetOneCreatedAndBindingReq) (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	condition := bson.M{}
	condition["service_code"] = data.ServiceCode
	condition["status"] = constants.ServicelogCreated
	condition["updated_date"] = bson.M{"$lt": now}

	update := bson.M{"$set": bson.M{
		"status":       constants.ServicelogBinding,
		"uid":          data.UID,
		"token":        data.Token,
		"android_id":   data.DeviceID,
		"updated_date": now,
	}}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
