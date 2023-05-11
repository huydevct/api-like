package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllGroupReq : get all friend
type AllGroupReq struct {
	Name   string
	Offset primitive.ObjectID
	Token  string
	Status []constants.CommonStatus
	Limit  int
}

// GroupInfo : Thông tin nhóm
type GroupInfo struct {
	ID     primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Name   string                 `json:"name" bson:"name"`
	Token  string                 `json:"token" bson:"token"`
	Status constants.CommonStatus `json:"status" bson:"status"`
	//
	CreatedIP   string     `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
}

// IsExists ..
func (m GroupInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
