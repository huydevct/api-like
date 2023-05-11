package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllGiftUsedReq : get all gift used
type AllGiftUsedReq struct {
	Code   string
	UserID primitive.ObjectID
	Offset primitive.ObjectID
	Limit  int
}

// GiftCodeUsed : Lưu thông tin gift code đã sử dụng
type GiftCodeUsed struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UserID    *primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Username  string              `json:"username,omitempty" bson:"username,omitempty"`
	Fullname  string              `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Code      string              `json:"code,omitempty" bson:"code,omitempty"`
	Value     int                 `json:"value,omitempty" bson:"value,omitempty"` // Giá trị gói
	CreatedAt int                 `json:"created_at,omitempty" bson:"created_at,omitempty"`
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
func (m GiftCodeUsed) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
