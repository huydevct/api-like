package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllGroupUIDReq : get all friend
type AllGroupUIDReq struct {
	UID     string
	GroupID primitive.ObjectID
	Token   string
	Status  []constants.CommonStatus
	Offset  primitive.ObjectID
	Limit   int
}

// GroupUIDInfo : Thông tin uid
type GroupUIDInfo struct {
	ID      primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UID     string                 `json:"uid" bson:"uid"`
	GroupID primitive.ObjectID     `json:"group_id" bson:"group_id"`
	Token   string                 `json:"token" bson:"token"`
	Status  constants.CommonStatus `json:"status" bson:"status"`
	//
	CreatedIP   string     `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	//
}

// IsExists ..
func (m GroupUIDInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
