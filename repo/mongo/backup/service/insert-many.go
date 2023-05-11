package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"app/model"
)

// InsertMany ..
func (r Repo) InsertMany(data []model.Service) (services []model.Service, err error) {

	collection := r.Session.GetCollection(r.Collection)
	ctx := r.GetContext()

	// Cạp nhật Data
	insertData := make([]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		data[i].ID = primitive.NewObjectID()
		//
		data[i].CreatedIP = data[i].UpdatedIP
		data[i].CreatedUser = data[i].UpdatedUser
		data[i].CreatedEmployee = data[i].UpdatedEmployee
		data[i].CreatedDate = data[i].UpdatedDate
		data[i].CreatedSource = data[i].UpdatedSource
		//
		data[i].CreatedAt = int(data[i].UpdatedDate.UnixNano() / int64(time.Millisecond))
		data[i].UpdateTime = int(data[i].UpdatedDate.UnixNano() / int64(time.Millisecond))

		insertData[i] = data[i]
		services = append(services, data[i])
	}

	_, err = collection.InsertMany(ctx, insertData, options.InsertMany().SetOrdered(false))
	if err != nil {
		return
	}
	return
}
