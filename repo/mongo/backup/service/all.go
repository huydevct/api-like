package service

import (
	"app/constants"
	"app/model"
	"app/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All : Lấy danh sách clone
func (r Repo) All(serviceType string, limit int) (results []model.ActiveService, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set option
	option := options.Find()
	option.SetSort(bson.M{"_id": 1})
	option.SetLimit(int64(limit))
	// set condition
	condition := bson.M{}
	if serviceType == "viplikeService" {
		startDate := utils.GetStartDateVNTime()
		endDate := utils.GetEndDateVNTime()
		condition["start_date"] = bson.M{"$gte": startDate, "$lte": endDate}
	}
	condition["status"] = constants.ServiceActive
	// condition["type"] = serviceType
	cur, err := collection.Find(ctx, condition, option)
	if err != nil {
		return
	}
	fmt.Println(condition, limit)
	// Close the cursor once finished
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		temp := model.ActiveService{}
		err := cur.Decode(&temp)
		if err != nil {
			return nil, err
		}

		results = append(results, temp)
	}

	return
}

// Search :
// . Sort
// . Filter
func (r Repo) Search(serviceType string, Page int, Limit int) (results []model.ActiveService, err error) {
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
	condition["action_type"] = serviceType
	condition["number_rest"] = bson.M{"$gte": 0}
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

func (r Repo) SearchByType(serviceType string, Page int, Limit int) (results []model.ActiveService, err error) {
	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	// set option
	option := options.Find()

	skip := (Page - 1) * Limit

	option.SetSkip(int64(skip))
	option.SetLimit(int64(Limit))

	// set condition
	condition := bson.M{}
	if serviceType == "viplikeService" {
		startDate := utils.GetStartDateVNTime()
		endDate := utils.GetEndDateVNTime()
		condition["start_date"] = bson.M{"$gte": startDate, "$lte": endDate}
	}
	condition["status"] = constants.ServiceActive
	condition["type"] = serviceType
	condition["number_rest"] = bson.M{"$gte": 0}
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
func (r Repo) Total(serviceType string) (total int, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	// set filter
	// set condition
	condition := bson.M{}
	condition["status"] = constants.ServiceActive
	condition["action_type"] = serviceType
	condition["number_rest"] = bson.M{"$gte": 0}

	count, err := collection.CountDocuments(ctx, condition)
	if err != nil {
		return
	}

	total = int(count)

	return
}

// Search :
// . Sort
// . Filter
func (r Repo) SearchType(serviceType string, Page int, Limit int) (results []model.ActiveService, err error) {
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
	if serviceType == "facebook_likepage_medium" {
		condition["type"] = bson.M{"$in": []string{"facebook_likepage_medium", "likepage_warranty"}}
	}
	if serviceType == "facebook_follow_medium" {
		condition["type"] = bson.M{"$in": []string{"facebook_follow_medium", "follow_warranty"}}
	}
	condition["number_rest"] = bson.M{"$gte": 0}
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
func (r Repo) TotalType(serviceType string) (total int, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	// set filter
	// set condition
	condition := bson.M{}
	condition["status"] = constants.ServiceActive
	condition["type"] = serviceType
	if serviceType == "viplikeService" {
		startDate := utils.GetStartDateVNTime()
		endDate := utils.GetEndDateVNTime()
		condition["start_date"] = bson.M{"$gte": startDate, "$lte": endDate}
	}
	if serviceType == "facebook_likepage_medium" {
		condition["type"] = bson.M{"$in": []string{"facebook_likepage_medium", "likepage_warranty"}}
	}
	if serviceType == "facebook_follow_medium" {
		condition["type"] = bson.M{"$in": []string{"facebook_follow_medium", "follow_warranty"}}
	}
	condition["number_rest"] = bson.M{"$gte": 0}

	count, err := collection.CountDocuments(ctx, condition)
	if err != nil {
		return
	}

	total = int(count)

	return
}

// AllViplike : Lấy tổng
func (r Repo) AllViplike() (total int, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	// set filter
	// set condition
	condition := bson.M{}
	condition["status"] = constants.ServiceActive
	condition["type"] = "viplikeService"
	startDate := utils.GetStartDateVNTime()
	condition["created_date"] = bson.M{"$gte": startDate}
	count, err := collection.CountDocuments(ctx, condition)
	if err != nil {
		return
	}

	total = int(count)

	return
}

// Search :
// . Sort
// . Filter
func (r Repo) SearchViplike(Page int, Limit int) (results []model.ActiveService, err error) {
	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	// set option
	option := options.Find()

	skip := (Page - 1) * Limit

	option.SetSkip(int64(skip))
	option.SetLimit(int64(Limit))
	condition := bson.M{}
	// set condition
	condition["status"] = constants.ServiceActive
	condition["type"] = "viplikeService"
	startDate := utils.GetStartDateVNTime()
	condition["created_date"] = bson.M{"$gte": startDate}
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
