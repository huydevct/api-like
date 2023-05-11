package service

import (
	"app/constants"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateServiceReq ..
type UpdateServiceReq struct {
	ServiceCode   string
	NumberSuccess int
	NumberRest    int
	IsSuccess     bool
}

// FindOneAndUpdateNumberSuccess : Cập nhật số luợng thành công theo service Code
func (r Repo) FindOneAndUpdateNumberSuccess(data UpdateServiceReq) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["service_code"] = data.ServiceCode

	// set update
	update := bson.M{}
	update["number_success"] = data.NumberSuccess
	update["number_success_int"] = data.NumberSuccess
	update["number_rest"] = data.NumberRest
	update["updated_date"] = now
	update["updateTime"] = int(now.UnixNano() / int64(time.Millisecond))

	if data.IsSuccess {
		update["status"] = constants.ServiceSuccess
	}

	updates := bson.M{
		"$set": update,
	}

	_, err = collection.UpdateOne(ctx,
		condition,
		updates)
	if err != nil {
		return
	}

	return
}
