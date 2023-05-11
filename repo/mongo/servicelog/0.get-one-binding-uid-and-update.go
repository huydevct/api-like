package servicelog

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByBindingUIDAndUpdate : Lấy 1 phần tử theo uid update binding -> getting
func (r Repo) GetOneByBindingUIDAndUpdate(cloneUID, serviceType string, appName constants.AppName) (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	condition := bson.M{}
	condition["uid"] = cloneUID
	condition["type"] = serviceType
	condition["appname"] = appName
	condition["status"] = constants.ServicelogBinding

	update := bson.M{"$set": bson.M{
		"status":       constants.ServicelogGetting,
		"updated_date": now,
	}}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
