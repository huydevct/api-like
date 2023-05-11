package clone

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchClonesByDeviceReq struct {
	DeviceSystem string
	DeviceID     primitive.ObjectID
	Limit        int
	Page         int
	OrderColumn  string
	OrderType    string
}

// GetAllCloneByDevice : Lấy danh sách clone theo device
func (r Repo) AllCloneByDevice(searchReq *SearchClonesByDeviceReq) (results []model.CloneInfo, err error) {
	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set sort
	option := options.Find()
	if searchReq.OrderColumn != "" {
		OrderType := -1
		if searchReq.OrderType == "asc" {
			OrderType = 1
		}
		option.SetSort(bson.M{searchReq.OrderColumn: OrderType})
	} else {
		option.SetSort(bson.M{"_id": -1})
	}

	// set page, limit
	if searchReq.Limit != -1 {
		if searchReq.Limit > 0 && searchReq.Limit <= 100 {
			option.SetLimit(int64(searchReq.Limit))
		} else {
			option.SetLimit(100)
		}
	}

	if searchReq.Page >= 0 {
		if searchReq.Page == 0 {
			searchReq.Page = 1
		}
		skip := (searchReq.Page - 1) * searchReq.Limit

		option.SetSkip(int64(skip))
	}

	// set filter
	filter := bson.M{}
	andFilter := []bson.M{}

	if !searchReq.DeviceID.IsZero() {
		deviceFilter := bson.M{}
		if searchReq.DeviceSystem != "" {
			switch searchReq.DeviceSystem {
			case "f_care":
				deviceFilter = bson.M{"chorm_device_id": searchReq.DeviceID}
			case "f_system":
				deviceFilter = bson.M{"ld_device_id": searchReq.DeviceID}
			case "f_android_webview":
				deviceFilter = bson.M{"f_view_device_id": searchReq.DeviceID}
			case "f_android":
				deviceFilter = bson.M{"android_device_id": searchReq.DeviceID}
			case "f_ios_webview":
				deviceFilter = bson.M{"f_ios_device_id": searchReq.DeviceID}
			case "f_ios":
				deviceFilter = bson.M{"iosdeviceid": searchReq.DeviceID}
			}
		} else {
			deviceFilter = bson.M{
				"$or": []bson.M{
					{"chorm_device_id": searchReq.DeviceID},
					{"ld_device_id": searchReq.DeviceID},
					{"f_view_device_id": searchReq.DeviceID},
					{"android_device_id": searchReq.DeviceID},
					{"f_ios_device_id": searchReq.DeviceID},
					{"iosdeviceid": searchReq.DeviceID},
				},
			}
		}
		andFilter = append(andFilter, deviceFilter)
	}

	aliveStatusFilter := bson.M{"alive_status": bson.M{"$ne": constants.CloneDelete}}
	andFilter = append(andFilter, aliveStatusFilter)

	filter["$and"] = andFilter

	cur, err := collection.Find(
		ctx,
		filter,
		option,
	)
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

		results = append(results, temp)
	}
	return
}

// All : Lấy danh sách clone mà device có
func (r Repo) TotalAllCloneByDevice(searchReq *SearchClonesByDeviceReq) (total int, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set filter
	filter := bson.M{}
	andFilter := []bson.M{}

	if !searchReq.DeviceID.IsZero() {
		deviceFilter := bson.M{
			"$or": []bson.M{
				{"android_device_id": searchReq.DeviceID},
				{"chorm_device_id": searchReq.DeviceID},
				{"f_ios_device_id": searchReq.DeviceID},
				{"iosdeviceid": searchReq.DeviceID},
				{"chorm_device_id": searchReq.DeviceID},
				{"ld_device_id": searchReq.DeviceID},
			},
		}
		andFilter = append(andFilter, deviceFilter)
	}

	aliveStatusFilter := bson.M{"alive_status": bson.M{"$ne": constants.CloneDelete}}
	andFilter = append(andFilter, aliveStatusFilter)

	filter["$and"] = andFilter

	//Đếm số lượng clone mà device có
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	total = int(count)

	return
}
