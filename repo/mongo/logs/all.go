package logs

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách clone
func (r Repo) All(req model.AllCloneRegReq) (results []model.CloneInfo, err error) {

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
	// set page
	if req.Page >= 0 {
		if req.Page == 0 {
			req.Page = 1
		}
		skip := (req.Page - 1) * req.Limit

		option.SetSkip(int64(skip))
	}

	// set filter
	filter := bson.M{}
	if len(req.AliveStatus) > 0 {
		filter["alive_status"] = bson.M{"$in": req.AliveStatus}
	}
	if req.StartDate != nil && req.EndDate == nil {
		filter["created_date"] = bson.M{"$gte": req.StartDate}
	} else if req.StartDate == nil && req.EndDate != nil {
		filter["created_date"] = bson.M{"$lte": req.EndDate}
	} else if req.StartDate != nil && req.EndDate != nil {
		filter["created_date"] = bson.M{"$gte": req.StartDate, "$lte": req.EndDate}
	}
	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.CloneInfo{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
