package servicelog

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneExpiredAndUpdate :
// 1. Tìm 1 phần tử có trạng thái "Binding", "Updating", "Getting" và updated_date + 4h < now
func (r Repo) GetOneExpiredAndUpdate() (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()
	expiredDate := now.Add(-1 * time.Hour)

	condition := bson.M{}
	condition["status"] = bson.M{"$in": []constants.ServiceLogStatus{
		constants.ServicelogBinding,
		constants.ServicelogUpdating,
		constants.ServicelogGetting,
	}}
	condition["updated_date"] = bson.M{"$lt": expiredDate}

	update := bson.M{"$set": bson.M{
		"status":       constants.ServicelogCreated,
		"uid":          "",
		"android_id":   "",
		"token":        "",
		"updated_date": now,
	}}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
