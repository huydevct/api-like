package model

import (
	"time"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllGiftReq : get all gift
type AllGiftReq struct {
	Status []constants.GiftStatus
	Code   string
	Offset primitive.ObjectID
	Limit  int
}

// GiftCode : Lưu thông tin gift code
type GiftCode struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Code      string               `json:"code,omitempty" bson:"code,omitempty"`
	Value     int                  `json:"value,omitempty" bson:"value,omitempty"` // Giá trị gói
	Status    constants.GiftStatus `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt int                  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
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
}

// IsExists ..
func (m GiftCode) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
