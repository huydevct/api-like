package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActiveClone : Lưu danh sách các clone đang active
type ActiveClone struct {
	ID        primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UID       string                 `json:"uid" bson:"uid"`         // (unit)
	AppName   constants.AppName      `json:"appname" bson:"appname"` // appName: facebook, instagram, tiktok,..
	Country   constants.CloneCountry `json:"country" bson:"country"` // type: vn, indo, en,..
	Token     string                 `json:"token" bson:"token"`
	ActionAt  *time.Time             `json:"action_at,omitempty" bson:"action_at,omitempty"`   // Lưu thời gian hoạt đông của clone
	ExpiredAt *time.Time             `json:"expired_at,omitempty" bson:"expired_at,omitempty"` // 8h kể từ khi submit
}

// IsExists ..
func (m ActiveClone) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
