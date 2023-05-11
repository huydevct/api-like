package groupuid

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách group uid
func (r Repo) All(req model.AllGroupUIDReq) (results []model.GroupUIDInfo, err error) {

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

	if req.GroupID.Hex() != "" {
		filter["group_id"] = req.GroupID
	}

	if len(req.Status) > 0 {
		filter["status"] = bson.M{"$in": req.Status}
	}
	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.GroupUIDInfo{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
