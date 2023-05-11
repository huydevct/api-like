package token

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson"
)

// ClearAll ..
func (r Repo) ClearAll(userID primitive.ObjectID, userType string) (matchCount int, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	condition := bson.M{
		"user_id": userID,
		"type":    userType,
		"status": bson.M{
			"$in": []constants.CommonStatus{
				constants.Active,
			},
		},
	}

	updateReq := bson.M{"status": constants.Delete}

	// Update data
	now := time.Now()
	updateReq["updated_date"] = &now
	updateResult, err := collection.UpdateMany(ctx, condition, bson.M{"$set": updateReq})
	if err != nil {
		return
	}

	matchCount = int(updateResult.MatchedCount)
	return
}

// ClearToken ..
func (r Repo) ClearToken(token, userType string) (matchCount int, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)

	condition := bson.M{
		"token": token,
		"type":  userType,
		"status": bson.M{
			"$in": []constants.CommonStatus{
				constants.Active,
			},
		},
	}

	updateReq := bson.M{"status": constants.Delete}

	// Update data
	now := time.Now()
	updateReq["updated_date"] = &now
	updateResult, err := collection.UpdateOne(ctx, condition, bson.M{"$set": updateReq})
	if err != nil {
		return
	}

	matchCount = int(updateResult.MatchedCount)
	return
}
