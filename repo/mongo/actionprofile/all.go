package actionprofile

import (
	"app/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách action profile
func (r Repo) All(req model.AllActionProfileReq) (results []model.ActionProfile, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"_id": -1})

	if req.Limit != -1 {
		if req.Limit > 0 && req.Limit <= 100 {
			option.SetLimit(int64(req.Limit))
		} else {
			option.SetLimit(100)
		}
	}

	// set filter
	filter := bson.M{}
	if !req.Offset.IsZero() {
		filter["_id"] = bson.M{"$lt": req.Offset}
	}

	// text search
	if req.Name != "" {
		filter["$text"] = bson.M{"$search": req.Name}
	}

	if req.Token != "" && req.Template != 1 {
		filter["token"] = req.Token
	}

	if req.Template > 0 {
		filter["template"] = req.Template
	}

	if req.AppName != "" {
		filter["appname"] = req.AppName
	}
	if len(req.Status) > 0 {
		filter["status"] = bson.M{"$in": req.Status}
	}
	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}
	fmt.Println(filter)
	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.ActionProfile{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
