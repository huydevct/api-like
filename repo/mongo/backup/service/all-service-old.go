package service

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Search :
// . Sort
// . Filter
func (r Repo) SearchAll(serviceType string, Page int, Limit int) (results []model.ActiveService, err error) {
	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	// set option
	option := options.Find()

	skip := (Page - 1) * Limit

	option.SetSkip(int64(skip))
	option.SetLimit(int64(Limit))

	// set condition
	condition := bson.M{}
	condition["status"] = constants.ServiceActive
	condition["type"] = serviceType
	cur, err := collection.Find(ctx, condition, option)
	if err != nil {
		return
	}
	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp model.ActiveService
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}
		results = append(results, temp)
	}

	return
}

// Total : Lấy tổng
func (r Repo) TotalTypeAll(serviceType string) (total int, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	// set filter
	// set condition
	condition := bson.M{}
	condition["status"] = constants.ServiceActive
	condition["type"] = serviceType
	count, err := collection.CountDocuments(ctx, condition)
	if err != nil {
		return
	}

	total = int(count)

	return
}
