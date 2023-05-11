package service

import (
	"app/constants"
	"app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateServiceReq ..
type UpdateServiceReq struct {
	ServiceCode   string
	NumberSuccess int
	NumberRest    int
	IsSuccess     bool
}

// IncNumberSuccess : Tăng số luợng success
func (r Repo) IncNumberSuccess(serviceCode string) (result model.ActiveService, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["service_code"] = serviceCode
	condition["status"] = constants.ServiceActive

	// set update
	update := bson.M{
		"$set": bson.M{
			"updated_date": now,
			"updated_by":   "api8",
		},
		"$inc": bson.M{
			"number_success": 1,
		},
	}

	myResult := collection.FindOneAndUpdate(ctx,
		condition,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	if myResult.Err() != nil {
		err = myResult.Err()
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}

	err = myResult.Decode(&result)
	if err != nil {
		return
	}

	return
}
