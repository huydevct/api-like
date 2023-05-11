package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ClearNumberReportByServiceCode : CẬp nhật số luợng report theo serviceCode
func (r Repo) ClearNumberReportByServiceCode(serviceCode string) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	// override data
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["service_code"] = serviceCode

	// set update
	update := bson.M{
		"$set": bson.M{
			"number_report": 0,
			"updated_date":  now,
			"resume_date":   now,
			"updated_by":    "",
		},
	}

	// set update
	_, err = collection.UpdateOne(ctx, condition, update)
	if err != nil {
		return
	}

	return
}
