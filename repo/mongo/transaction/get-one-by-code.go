package transaction

import (
	"app/model"
	"app/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByCode ..
func (r Repo) GetOneByCode(code string) (result model.Transaction, err error) {

	ctx := context.Background()

	//
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"code": code}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}
	// TODO : cập nhật giá tri
	{
		if result.IsExists() && result.ValueInt == 0 {
			result.ValueInt = utils.ConvertToInt(result.Value)
		}
	}

	return
}
