package servicelog

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetOneCreatedAndFinishedReq ..
type GetOneCreatedAndFinishedReq struct {
	ServiceCode    string
	UID            string
	ViplikePostID  string
	SharePostID    primitive.ObjectID
	ShareCommentID primitive.ObjectID
	DeviceID       string
	Token          string
}

// GetOneCreatedAndFinished : Tìm 1 gói created và update thành finished
func (r Repo) GetOneCreatedAndFinished(req GetOneCreatedAndFinishedReq) (result model.ServiceLog, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()

	// set condition
	condition := bson.M{}
	condition["service_code"] = req.ServiceCode
	condition["status"] = constants.ServicelogCreated

	// set update
	update := bson.M{
		"$set": bson.M{
			"status":           constants.ServicelogFinished,
			"uid":              req.UID,
			"device_id":        req.DeviceID,
			"token":            req.Token,
			"share_comment_id": req.ShareCommentID,
			"share_post_id":    req.SharePostID,
			"viplike_post_id":  req.ViplikePostID,
			"updated_date":     now,
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
