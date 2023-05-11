package walletlog

import (
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách wallet log
func (r Repo) All(req model.AllWalletLogReq) (results []model.WalletLog, err error) {

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
	if req.UserToken != "" {
		filter["token"] = req.UserToken
	}
	if len(req.Types) > 0 {
		filter["type"] = bson.M{"$in": req.Types}
	}
	if req.FromTime > 0 && req.ToTime == 0 {
		fromTime := time.Unix(0, int64(req.FromTime)*int64(time.Millisecond))
		filter["created_date"] = bson.M{"$gt": fromTime}
	} else if req.FromTime == 0 && req.ToTime > 0 {
		toTime := time.Unix(0, int64(req.ToTime)*int64(time.Millisecond))
		filter["created_date"] = bson.M{"$lt": toTime}
	} else if req.FromTime > 0 && req.ToTime > 0 {
		fromTime := time.Unix(0, int64(req.FromTime)*int64(time.Millisecond))
		toTime := time.Unix(0, int64(req.ToTime)*int64(time.Millisecond))
		filter["created_date"] = bson.M{"$gt": fromTime, "$lt": toTime}
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.WalletLog{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
