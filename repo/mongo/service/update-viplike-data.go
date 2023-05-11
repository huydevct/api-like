package service

import (
	"app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateViplikeDataReq ..
type UpdateViplikeDataReq struct {
	ServiceID      primitive.ObjectID
	P1             bool
	P2             bool
	P3             bool
	P4             bool
	P5             bool
	ViplikeDataOld []model.ViplikeItemOld
}

// UpdateViplikeData : Cập nhật data
func (r Repo) UpdateViplikeData(req UpdateViplikeDataReq) (result model.Service, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["_id"] = req.ServiceID

	// set update
	updates := bson.M{
		"$set": bson.M{
			"updated_by":   "api8",
			"updateTime":   int(now.UnixNano() / int64(time.Millisecond)),
			"updated_date": now,
		},
	}
	inc := bson.M{}
	if req.P1 {
		inc["data_viplike.p1"] = 1
	}
	if req.P2 {
		inc["data_viplike.p2"] = 1
	}
	if req.P3 {
		inc["data_viplike.p3"] = 1
	}
	if req.P4 {
		inc["data_viplike.p4"] = 1
	}
	if req.P5 {
		inc["data_viplike.p5"] = 1
	}

	// Chỉ cập nhật nếu có thay đổi bài post
	if len(inc) > 0 {
		updates["$inc"] = inc

		myResult := collection.FindOneAndUpdate(ctx,
			condition,
			updates,
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
	}

	return
}

// UpdateViplikeDataOld : Cập nhật data old
func (r Repo) UpdateViplikeDataOld(req UpdateViplikeDataReq) (result model.Service, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["_id"] = req.ServiceID

	// set update
	update := bson.M{
		"$set": bson.M{
			"data":         req.ViplikeDataOld,
			"updated_date": now,
			"updateTime":   int(now.UnixNano() / int64(time.Millisecond)),
			"updated_by":   "api8",
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
