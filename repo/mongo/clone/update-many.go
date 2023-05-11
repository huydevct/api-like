package clone

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Update ..
func (r Repo) UpdateMany() (err error) {

	ctx := r.GetContext()
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"alive_status": "getting"}
	_, err = collection.UpdateMany(ctx,
		condition,
		bson.M{"$set": bson.M{"alive_status": "live"}})
	if err != nil {
		return
	}
	return
}
