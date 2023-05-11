package pack

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách packages theo status, offset, limit
func (r Repo) All(data model.AllPackageReq) (results []model.Package, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"_id": -1})

	if data.Limit != -1 {
		if data.Limit > 0 && data.Limit <= 100 {
			option.SetLimit(int64(data.Limit))
		} else {
			option.SetLimit(100)
		}
	}
	// set filter
	filter := bson.M{}
	if !data.Offset.IsZero() {
		filter["_id"] = bson.M{"$lt": data.Offset}
	}
	if len(data.Status) > 0 {
		filter["status"] = bson.M{"$in": data.Status}
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.Package{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
