package servicelog

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách servicelog
func (r Repo) All(req model.AllServicelogReq) (results []model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetLimit(int64(req.Limit))

	// set filter
	filter := bson.M{}
	if req.UID != "" {
		filter["uid"] = req.UID
	}
	if req.Type != "" {
		filter["type"] = req.Type
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.ServiceLog{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
