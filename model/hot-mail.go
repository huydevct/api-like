package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllHotMailReq : get all hotmail
type AllHotMailReq struct {
	Name   string
	Offset primitive.ObjectID
	Token  string
	Status []constants.CommonStatus
	Limit  int
}

// HotMail ..
type HotMail struct {
	ID          primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	Email       string                  `json:"email,omitempty" bson:"email,omitempty"`
	PassWord    string                  `json:"password,omitempty" bson:"password,omitempty"`
	Device      string                  `json:"device,omitempty" bson:"device,omitempty"`
	Status      constants.HotMailStatus `json:"status,omitempty" bson:"status,omitempty"`
	CreatedDate *time.Time              `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate *time.Time              `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
}

// IsExists struct
func (m HotMail) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
