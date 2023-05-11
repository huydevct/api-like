package gift

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApplyWithSessionCtx ..
func (r Repo) ApplyWithSessionCtx(sessionCtx context.Context, code string) (result model.GiftCode, err error) {

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"code": code, "status": constants.GiftActive}
	updateReq := bson.M{"$set": bson.M{"status": constants.GiftUsed, "updated_date": time.Now()}}

	err = collection.FindOneAndUpdate(sessionCtx, condition, updateReq).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
