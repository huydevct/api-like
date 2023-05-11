package frienduid

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllByFriendIDUID : ..
func (r Repo) AllByFriendIDUID(friendID primitive.ObjectID, lastFriendUID primitive.ObjectID, quantity int) (results []model.FriendUIDInfo, uids []string, ids []primitive.ObjectID, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetLimit(int64(quantity))

	// set filter
	filter := bson.M{}
	filter["friend_id"] = friendID
	if !lastFriendUID.IsZero() {
		filter["_id"] = bson.M{"$gt": lastFriendUID}
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.FriendUIDInfo{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, nil, nil, err
		}

		uids = append(uids, temp.UID)
		ids = append(ids, temp.ID)
		results = append(results, temp)
	}

	return
}
