package groupuid

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertByUID : upsert by uid
func (r Repo) UpsertByUID(data model.GroupUIDInfo) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	_, err = collection.UpdateOne(ctx,
		bson.M{"uid": data.UID, "group_id": data.GroupID},
		bson.M{"$set": data},
		options.Update().SetUpsert(true))

	return
}
