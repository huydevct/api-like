package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllFriendReq : get all friend
type AllFriendReq struct {
	Name   string
	Offset primitive.ObjectID
	Token  string
	Status []constants.CommonStatus
	Limit  int
}

// FriendInfo : Thông tin bạn bè
type FriendInfo struct {
	ID     primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Name   string                 `json:"name" bson:"name"`
	Token  string                 `json:"token" bson:"token"`
	Status constants.CommonStatus `json:"status" bson:"status"`
	//
	UpdatedIP       string              `json:"updated_ip,omitempty" bson:"updated_ip,omitempty"`
	UpdatedEmployee *primitive.ObjectID `json:"updated_employee,omitempty" bson:"updated_employee,omitempty"`
	UpdatedUser     *primitive.ObjectID `json:"updated_user,omitempty" bson:"updated_user,omitempty"`
	UpdatedSource   string              `json:"updated_source,omitempty" bson:"updated_source,omitempty"`
	UpdatedDate     *time.Time          `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	//
	CreatedIP       string              `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedEmployee *primitive.ObjectID `json:"created_employee,omitempty" bson:"created_employee,omitempty"`
	CreatedUser     *primitive.ObjectID `json:"created_user,omitempty" bson:"created_user,omitempty"`
	CreatedSource   string              `json:"created_source,omitempty" bson:"created_source,omitempty"`
	CreatedDate     *time.Time          `json:"created_date,omitempty" bson:"created_date,omitempty"`
	//
	LastChangePasswordDate *time.Time `json:"-" bson:"last_change_password_date,omitempty"`
}

// IsExists ..
func (m FriendInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
