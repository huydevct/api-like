package user

import (
	"app/constants"
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByAgentCode ..
func (r Repo) GetOneByAgentCode(AgentCode string) (result model.UserInfo, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"agent_code": AgentCode, "status": constants.Active}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}

// GetOneByInviteCode ..
func (r Repo) GetOneByInviteCode(InviteCode string) (result model.UserInfo, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	condition := bson.M{"user_invite_code": InviteCode, "status": constants.Active}
	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
