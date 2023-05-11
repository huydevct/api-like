package service

import (
	"app/constants"
	"app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByUID ..
func (r Repo) GetOneByUID(uid string, token string) (result model.Service, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"fanpage_id": uid, "status": constants.ServiceSuccess, "token": token}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneByServiceActive ..
func (r Repo) GetOneByServiceActive(service_code string, token string) (result model.Service, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	filter := bson.M{}
	filter["service_code_old"] = service_code
	orQuery := []bson.M{}
	orQuery = append(orQuery, bson.M{"status": constants.ServiceActive, "token": token}, bson.M{"user_token": token, "status": bson.M{"$in": []constants.ServiceStatus{constants.ServiceActive, constants.ServicePendingWarranty}}})
	filter["$or"] = orQuery
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneByServiceCode ..
func (r Repo) GetOneByServiceCode(service_code string, token string) (result model.Service, err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	filter := bson.M{}
	filter["service_code"] = service_code
	filter["token"] = token
	filter["status"] = constants.ServiceSuccess
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
