package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllCloneReq : get all clone
type AllCloneReq struct {
	UID         string
	Token       string
	Email       string
	DeviceID    string
	AliveStatus []constants.AliveStatus
	Offset      primitive.ObjectID
	Limit       int
	Page        int
}

// AllCloneRegReq : get all clone
type AllCloneRegReq struct {
	AliveStatus []constants.AliveStatus
	Limit       int
	Page        int
	StartDate   *time.Time
	EndDate     *time.Time
}

// CloneInfo : Thông tin clone
type CloneInfo struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Email       string                 `json:"email" bson:"email"`
	UID         string                 `json:"uid" bson:"uid"`
	IP          string                 `json:"ip" bson:"ip"`
	PCName      string                 `json:"pc_name" bson:"pc_name"`
	Token       string                 `json:"token" bson:"token"`
	DeviceID    string                 `json:"device_id" bson:"device_id"` // Lưu để kiểm tra xem clone này đang thuộc thiết bị nào
	AliveStatus constants.AliveStatus  `json:"alive_status" bson:"alive_status"`
	Password    string                 `json:"password" bson:"password"`
	Secretkey   string                 `json:"secretkey" bson:"secretkey"`
	Cookie      string                 `json:"cookie" bson:"cookie"`
	Language    string                 `json:"language" bson:"language"` // Ngôn ngữ
	Country     constants.CloneCountry `json:"country" bson:"country"`   // kiểm tra xem clone của nước nào
	AppName     string                 `json:"appname" bson:"appname"`   // loại clone : facebook, instagram, tiktok...

	// action
	ActionProfileID *primitive.ObjectID `json:"action_profile_id,omitempty" bson:"action_profile_id,omitempty"`

	// Clone Info
	Name        string `json:"name" bson:"name,omitempty"`
	Birthday    string `json:"birthday" bson:"birthday,omitempty"`
	PhoneNumber string `json:"phone_number" bson:"phone_number,omitempty"`
	Follow      int    `json:"follow" bson:"follow,omitempty"`
	Sex         string `json:"sex" bson:"sex,omitempty"`
	Friend      int    `json:"friend" bson:"friend,omitempty"`
	// trace

	CreatedCloneDate string     `json:"created_clone_date,omitempty" bson:"created_clone_date,omitempty"`
	CreatedDate      *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate      *time.Time `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
}

// IsExists ..
func (m CloneInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// SearchClone : Thông tin clone
type SearchClone struct {
	AliveStatus []constants.AliveStatus
	AppName     []constants.AppName
	DeviceID    string
	System      string
	DeviceObID  primitive.ObjectID
	UID         string
	CloneID     string
	Date        interface{}
	IsReg       *bool
	Token       string // Autofarmer token
	Page        int
	Limit       int
}
