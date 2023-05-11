package transaction

import (
	"app/model"
	"app/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách gói giao dịch theo status, userID, offset, limit
func (r Repo) All(req model.AllTransactionReq) (results []model.Transaction, err error) {

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
	if req.Code != "" {
		filter["code"] = req.Code
	}
	if req.UserToken != "" {
		filter["token"] = req.UserToken
	}
	if len(req.Status) > 0 {
		filter["status"] = bson.M{"$in": req.Status}
	}
	// text search
	if req.Username != "" {
		filter["$text"] = bson.M{"$search": req.Username}
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}

	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.Transaction{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		// TODO : cập nhật giá tri
		{
			if temp.IsExists() && temp.ValueInt == 0 {
				temp.ValueInt = utils.ConvertToInt(temp.Value)
			}
		}

		results = append(results, temp)
	}

	return
}
