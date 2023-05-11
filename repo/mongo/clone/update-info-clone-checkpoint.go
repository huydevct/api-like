package clone

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateCloneChekpointInfoReq
type UpdateCloneChekpointInfoReq struct {
	CloneID primitive.ObjectID
	IsLive  string
}

// UpdateCloneChekpointInfo ..
func (r Repo) UpdateCloneChekpointInfo(req UpdateCloneChekpointInfoReq) (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["_id"] = req.CloneID

	// set update
	update := bson.M{}
	update["updated_clone_checkpoint"] = time.Now()

	//========

	if req.IsLive == "true" {
		update["alive_status"] = constants.CloneLive
		update["android_id"] = ""
		update["updated_by"] = "reset_clone"
	}

	updates := bson.M{
		"$set": update,
	}

	myResult := collection.FindOneAndUpdate(ctx,
		condition,
		updates,
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	if myResult.Err() != nil {
		err = myResult.Err()
		return
	}

	err = myResult.Decode(&result)
	if err != nil {
		return
	}

	return
}
