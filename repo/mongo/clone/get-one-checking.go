package clone

import (
	"app/constants"
	"app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneChecking :
// 1. Tìm 1 phần tử có trạng thái "Checking" và updated_date + 24h < now
func (r Repo) GetOneChecking() (result model.CloneInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	now := time.Now()
	expiredDate := now.Add(-6 * time.Hour)

	condition := bson.M{}
	condition["alive_status"] = constants.CloneChecking
	condition["updated_date"] = bson.M{"$lt": expiredDate}

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
