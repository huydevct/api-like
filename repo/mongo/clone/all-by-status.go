package clone

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllReq : Lấy danh sách gói success
type AllCloneReq struct {
	AliveStatus constants.AliveStatus
	Offset      primitive.ObjectID
	Limit       int
}

// AllCloneByStatus :
func (r Repo) AllCloneByStatus(req AllCloneReq) (results []model.CloneInfo, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	// set time sau 72h
	now := time.Now()
	then := now.AddDate(0, 0, -3)

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
	if req.AliveStatus != "" {
		condition["alive_status"] = req.AliveStatus
	}
	condition["updated_date"] = bson.M{
		"$lte": then,
	}

	cur, err := collection.Find(ctx, condition, option)
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
