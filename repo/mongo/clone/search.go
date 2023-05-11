package clone

import (
	"app/model"
	"app/utils"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Search :
// . Sort
// . Filter
func (r Repo) Search(searchReq model.SearchClone) (results []model.CloneInfo, err error) {
	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	// set option
	option := options.Find()

	// set page, limit
	if searchReq.Limit == 0 || searchReq.Limit > 1000 {
		searchReq.Limit = 1000
	}
	if searchReq.Page == 0 {
		searchReq.Page = 1
	}
	skip := (searchReq.Page - 1) * searchReq.Limit

	option.SetSkip(int64(skip))
	option.SetLimit(int64(searchReq.Limit))

	// Sort default theo _id
	option.SetSort(bson.M{"_id": -1})

	// set filter
	filter := bson.M{}
	if searchReq.Token != "" {
		filter["token"] = searchReq.Token
	}
	if len(searchReq.AliveStatus) > 0 {
		filter["alive_status"] = bson.M{"$in": searchReq.AliveStatus}
	}
	if searchReq.Date != nil && searchReq.Date.(string) != "" {
		date, _ := utils.ConvertDateReturn(searchReq.Date.(string))
		date = date.Add(time.Hour * 7)
		StartDate, EndDate := utils.GetStartEndOfDate(date)
		fmt.Println(StartDate, EndDate)
		filter["updated_date"] = bson.M{"$gte": StartDate, "$lte": EndDate}
	}

	if len(searchReq.AppName) > 0 {
		filter["appname"] = bson.M{"$in": searchReq.AppName}
	}

	if searchReq.IsReg != nil {
		filter["is_reg"] = searchReq.IsReg
	}

	cur, err := collection.Find(ctx, filter, option)
	if err != nil {
		return
	}
	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp model.CloneInfo
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}
		temp.Password = ""
		results = append(results, temp)
	}

	return
}
