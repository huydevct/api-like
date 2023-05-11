package clone

import (
	"app/model"
	"app/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Total : Lấy tổng
func (r Repo) Total(searchReq model.SearchClone) (total int, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set filter
	filter := bson.M{}
	if searchReq.Token != "" {
		filter["token"] = searchReq.Token
	}
	if len(searchReq.AliveStatus) > 0 {
		filter["alive_status"] = bson.M{"$in": searchReq.AliveStatus}
	}
	if len(searchReq.AppName) > 0 {
		filter["appname"] = bson.M{"$in": searchReq.AppName}
	}

	if searchReq.Date != nil && searchReq.Date.(string) != "" {
		date, _ := utils.ConvertDateReturn(searchReq.Date.(string))
		date = date.Add(time.Hour * 7)
		StartDate, EndDate := utils.GetStartEndOfDate(date)
		filter["updated_date"] = bson.M{"$gte": StartDate, "$lte": EndDate}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	total = int(count)

	return
}
