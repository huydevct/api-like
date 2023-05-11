package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllFriendUIDReq : get all friend
type AllFriendUIDReq struct {
	UID      string
	FriendID primitive.ObjectID
	Token    string
	Status   []constants.CommonStatus
	Offset   primitive.ObjectID
	Limit    int
}

// FriendUIDInfo : Thông tin uid
type FriendUIDInfo struct {
	ID       primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UID      string                 `json:"uid" bson:"uid"`
	FriendID primitive.ObjectID     `json:"friend_id" bson:"friend_id"`
	Token    string                 `json:"token" bson:"token"`
	Status   constants.CommonStatus `json:"status" bson:"status"`

	//
	CreatedIP   string     `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	//
}

// IsExists ..
func (m FriendUIDInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
