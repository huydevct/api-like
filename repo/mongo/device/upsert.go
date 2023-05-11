package device

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Upsert : upsert thông tin device
// 2. DeviceID
func (r Repo) Upsert(data model.DeviceInfo) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	condition := bson.M{}
	condition["pc_name"] = data.PCName
	condition["status"] = bson.M{"$ne": constants.Delete}

	data.UpdatedDate = &now

	_, err = collection.UpdateOne(ctx,
		condition,
		bson.M{"$set": data},
		options.Update().SetUpsert(true))

	return
}
