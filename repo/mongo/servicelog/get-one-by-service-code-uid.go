package servicelog

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByServiceCodeUID :
func (r Repo) GetOneByServiceCodeUID(serviceCode, cloneUID string) (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	condition := bson.M{}
	condition["service_code"] = serviceCode
	condition["uid"] = cloneUID

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneByServiceCodeUIDViplikePostID :
func (r Repo) GetOneByServiceCodeUIDViplikePostID(serviceCode, cloneUID, viplikePostID string) (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	condition := bson.M{}
	condition["service_code"] = serviceCode
	condition["uid"] = cloneUID
	condition["viplike_post_id"] = viplikePostID

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
