package device

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Search :
// . Sort
// . Filter
func (r Repo) Search(searchReq model.SearchDevice) (results []model.DeviceInfo, err error) {
	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	// set option
	option := options.Find()

	// set page, limit
	if searchReq.Limit == 0 || searchReq.Limit > 500 {
		searchReq.Limit = 500
	}
	if searchReq.Page == 0 {
		searchReq.Page = 1
	}
	skip := (searchReq.Page - 1) * searchReq.Limit

	option.SetSkip(int64(skip))
	option.SetLimit(int64(searchReq.Limit))

	// Sort default theo _id
	option.SetSort(bson.M{"updated_date": -1})

	// set filter
	filter := bson.M{}
	if searchReq.Token != "" {
		filter["token"] = searchReq.Token
	}
	if len(searchReq.Status) > 0 {
		filter["status"] = bson.M{"$in": searchReq.Status}
	}
	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}
	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp model.DeviceInfo
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}
		results = append(results, temp)
	}

	return
}
