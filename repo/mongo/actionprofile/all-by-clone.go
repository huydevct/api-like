package actionprofile

import (
	"app/constants"
	"app/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllByClone : Lấy danh sách action profile gán cho clone
func (r Repo) AllByClone(req model.AllActionProfileReq) (results []model.ActionProfile, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"template": 1})

	// set filter
	filter := bson.M{}
	if !req.Offset.IsZero() {
		filter["_id"] = bson.M{"$lt": req.Offset}
	}

	// text search
	if req.Name != "" {
		filter["$text"] = bson.M{"$search": req.Name}
	}
	filter["$or"] = []interface{}{
		bson.D{{"token", req.Token}},
		bson.D{{"template", constants.TemplateActionAdmin}},
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
