package user

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách user theo status, offset, limit
func (r Repo) All(data model.AllUserReq) (results []model.UserInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"device_total": -1})
	option.SetLimit(100000)
	// if data.Limit != -1 {
	// 	if data.Limit > 0 && data.Limit <= 100 {
	// 		option.SetLimit(int64(data.Limit))
	// 	} else {
	// 		option.SetLimit(100000)
	// 	}
	// }
	// set filter
	filter := bson.M{}
	if !data.Offset.IsZero() {
		filter["_id"] = bson.M{"$lt": data.Offset}
	}
	if len(data.Status) > 0 {
		filter["status"] = bson.M{"$in": data.Status}
	}
	// text search
	if data.Username != "" {
		filter["username"] = bson.M{"$regex": data.Username}
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.UserInfo{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
