package model

import (
	"time"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token : token đăng nhập của user
type Token struct {
	ID                 primitive.ObjectID     `json:"-" bson:"_id,omitempty"`
	Token              string                 `json:"token" bson:"token"`
	UserToken          string                 `json:"user_token,omitempty" bson:"user_token,omitempty"`
	Status             constants.CommonStatus `json:"status,omitempty" bson:"status,omitempty"`
	CreatedDate        *time.Time             `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate        *time.Time             `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	ExpiredAfterSecond int                    `json:"-" bson:"-"`                                           // Không lưu field này vào DB
	ExpiredDate        *time.Time             `json:"expired_date,omitempty" bson:"expired_date,omitempty"` // ExpiredDate = thời gian tạo + time expired duration
	LastUsedDate       *time.Time             `json:"last_used_date,omitempty" bson:"last_used_date,omitempty"`
	// token info
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Type      string             `json:"type" bson:"type"` // employee, user
	DeviceID  string             `json:"device_id,omitempty" bson:"device_id,omitempty"`
	FcmID     string             `json:"fcm_id,omitempty" bson:"fcm_id,omitempty"` // Lưu Firebase ID của app, support gửi notify đến app
	UserAgent string             `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	RemoteIP  string             `json:"remote_ip,omitempty" bson:"remote_ip,omitempty"`
	Source    string             `json:"source,omitempty" bson:"source,omitempty"`
}

// IsExists ..
func (m *Token) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
