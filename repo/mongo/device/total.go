package device

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// Total : Lấy tổng
func (r Repo) Total(searchReq model.SearchDevice) (total int, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set filter
	filter := bson.M{}
	if searchReq.Token != "" {
		filter["token"] = searchReq.Token
	}
	if len(searchReq.Status) > 0 {
		filter["status"] = bson.M{"$in": searchReq.Status}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	total = int(count)

	return
}
