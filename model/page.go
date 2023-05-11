package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllPageReq : get all page
type AllPageReq struct {
	UID   string
	Token string

	Status []constants.CommonStatus

	Offset primitive.ObjectID
	Limit  int
}

// PageInfo : Thông tin page
type PageInfo struct {
	ID     primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UID    string                 `json:"uid" bson:"uid"`
	Token  string                 `json:"token" bson:"token"`
	Status constants.CommonStatus `json:"status" bson:"status"`
	//
	CreatedIP   string     `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
}

// IsExists ..
func (m PageInfo) IsExists() (ok bool) {

	if !m.ID.IsZero() {
		ok = true
	}
	return
}
