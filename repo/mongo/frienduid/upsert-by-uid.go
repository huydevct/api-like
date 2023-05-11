package frienduid

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertByUID : upsert by uid
func (r Repo) UpsertByUID(data model.FriendUIDInfo) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	_, err = collection.UpdateOne(ctx,
		bson.M{"uid": data.UID, "friend_id": data.FriendID},
		bson.M{"$set": data},
		options.Update().SetUpsert(true))

	return
}
