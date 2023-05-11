package servicelog

import (
	"context"
	"time"

	"app/constants"
	"app/model"

	"go.mongodb.org/mongo-driver/bson"
)

// Update ..
func (r Repo) Update(data model.ServiceLog) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	data.UpdatedDate = &now
	condition := bson.M{}
	condition["_id"] = data.ID
	condition["status"] = bson.M{"$nin": []constants.ServiceLogStatus{
		constants.ServicelogBinding,
		constants.ServicelogGetting,
		constants.ServicelogFinished,
	}}

	_, err = collection.UpdateOne(ctx,
		condition,
		bson.M{"$set": data})
	if err != nil {
		return
	}

	return
}
