package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeviceInfo ..
type DeviceInfo struct {
	ID          primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	Token       string                 `json:"token,omitempty" bson:"token,omitempty"`
	PCName      string                 `json:"pc_name,omitempty" bson:"pc_name,omitempty"`
	Status      constants.CommonStatus `json:"status,omitempty" bson:"status,omitempty"`
	CreatedDate *time.Time             `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate *time.Time             `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
}

// IsExists ..
func (m DeviceInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// AllDeviceReq : get all clone
type AllDeviceReq struct {
	Name     string
	Token    string
	DeviceID string
	Status   []constants.CommonStatus
	Offset   primitive.ObjectID
	Limit    int
	Page     int
}

// SearchDevice : Thông tin device
type SearchDevice struct {
	Token  string // Autofarmer token
	Status []constants.CommonStatus
	Page   int
	Limit  int
}
