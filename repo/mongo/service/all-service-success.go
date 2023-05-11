package service

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllReq : Lấy danh sách gói success
type AllServiceReq struct {
	Offset    primitive.ObjectID
	Type      string
	Limit     int
	StartDate int
}

// AllServiceToday :
func (r Repo) AllServiceSuccess(req AllServiceReq) (results []model.Service, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"_id": 1})

	if req.Limit != -1 {
		if req.Limit > 0 && req.Limit <= 100 {
			option.SetLimit(int64(req.Limit))
		} else {
			option.SetLimit(100)
		}
	}

	// set condition
	condition := bson.M{}
	condition["type"] = req.Type
	condition["token"] = bson.M{"$ne": "WZ6QT84EX5KAZVK9GZQUY7S34Z8GRWHK"}
	condition["TimeSuccess"] = bson.M{
		"$gte": req.StartDate,
	}
	if !req.Offset.IsZero() {
		condition["_id"] = bson.M{"$gt": req.Offset}
	}

	cur, err := collection.Find(ctx, condition, option)
	if err != nil {
		return
	}
	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.Service{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}
