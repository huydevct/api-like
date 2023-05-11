package clone

import (
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByIDAndUpdateGetActionAt : Lấy phần tử theo ID và update get action at
func (r Repo) FindOneAndUpdateDoResultActionAt(cloneID primitive.ObjectID) (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	condition := bson.M{}
	condition["_id"] = cloneID

	update := bson.M{"$set": bson.M{
		"do_result_at": now,
		"updated_date": now,
	}}

	err = collection.FindOneAndUpdate(ctx, condition, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
