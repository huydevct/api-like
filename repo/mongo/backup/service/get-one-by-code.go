package service

import (
	"app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByCode ..
func (r Repo) GetOneByCode(serviceCode string) (result model.Service, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"service_code": serviceCode}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
