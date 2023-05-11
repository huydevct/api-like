package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CloneActivityInfo : Thông tin clone
type CloneActivityInfo struct {
	ID              primitive.ObjectID            `json:"id" bson:"_id,omitempty"`
	RequestDeviceID string                        `json:"request_android_id" bson:"request_android_id"` // DeviceID thực hiện getDoAction
	RequestID       string                        `json:"request_id" bson:"request_id"`                 // ID của body request
	ServiceCode     string                        `json:"service_code" bson:"service_code"`
	UID             string                        `json:"uid,omitempty" bson:"uid,omitempty"`
	Email           string                        `json:"email,omitempty" bson:"email,omitempty"`
	Token           string                        `json:"token" bson:"token"`
	DeviceID        string                        `json:"android_id" bson:"android_id"` // Lưu để kiểm tra xem clone này đang thuộc thiết bị nào
	Appname         constants.AppName             `json:"appname" bson:"appname"`       // autofarmer, facebook, youtube,..
	Action          string                        `json:"action" bson:"action"`
	Status          constants.CloneActivityStatus `json:"status" bson:"status"`
	Money           int                           `json:"money" bson:"money"`
	//
	CreatedIP   string     `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate *time.Time `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
}

// CloneActivityPublish : chứa data publish lên queue upsert clone_activities
type CloneActivityPublish struct {
	ServiceCode string                 `json:"service_code"`
	UID         string                 `json:"uid"`
	Email       string                 `json:"email"`
	Data        map[string]interface{} `json:"data"`
}

// IsExists ..
func (m CloneActivityInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
