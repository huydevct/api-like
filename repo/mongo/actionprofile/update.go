package actionprofile

import (
	"context"
	"fmt"

	"app/model"

	"go.mongodb.org/mongo-driver/bson"
)

// Update ..
func (r Repo) Update(data model.ActionProfile) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"_id": data.ID}

	updateResult, err := collection.UpdateOne(ctx,
		condition,
		bson.M{"$set": data})
	if err != nil {
		return
	}
	if updateResult.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}
