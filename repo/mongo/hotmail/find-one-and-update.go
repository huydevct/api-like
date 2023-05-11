package hotmail

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOneAndUpdate : update hotmail
func (r Repo) FindOneAndUpdate(Token string) (result model.HotMail, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()
	queryOptions := options.FindOneAndUpdateOptions{}
	queryOptions.SetUpsert(false)
	update := bson.M{"$set": bson.M{"token": Token, "updated_date": &now, "status": constants.HotMailUsed}}
	err = collection.FindOneAndUpdate(ctx, bson.M{"status": constants.HotMailLive}, update, &queryOptions).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}
	return
}
