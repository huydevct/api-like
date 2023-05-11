package servicelog

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByServiceCode : ..
func (r Repo) GetOneByServiceCode(serviceCode string) (result model.ServiceLog, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["service_code"] = serviceCode

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
