package gift

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByCode ..
func (r Repo) GetOneByCode(code string) (result model.GiftCode, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{
		"code": code,
		"status": bson.M{
			"$in": []constants.GiftStatus{
				constants.GiftActive,
				constants.GiftWaiting,
			},
		}}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
