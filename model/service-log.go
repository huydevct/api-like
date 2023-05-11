package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceLog : Lưu thông tin service log
type ServiceLog struct {
	ID                 primitive.ObjectID         `json:"id" bson:"_id,omitempty"`
	ServiceCode        string                     `json:"service_code,omitempty" bson:"service_code,omitempty"`
	ViplikeServiceCode string                     `json:"viplike_service_code,omitempty" bson:"viplike_service_code,omitempty"` // Mã dịch vụ dùng cho viplikeService
	ViplikePostID      string                     `json:"viplike_post_id,omitempty" bson:"viplike_post_id,omitempty"`           // P1, P2, .., P5
	Status             constants.ServiceLogStatus `json:"status" bson:"status"`
	Kind               constants.ServiceKind      `json:"kind,omitempty" bson:"kind,omitempty"`
	Price              int                        `json:"price,omitempty" bson:"price,omitempty"`
	UID                string                     `json:"uid,omitempty" bson:"uid"`
	DeviceID           string                     `json:"device_id" bson:"device_id"`
	Token              string                     `json:"token,omitempty" bson:"token"`
	AutolikeToken      string                     `json:"autolike_token,omitempty" bson:"autolike_token"`
	// Data service
	Type          constants.ServiceType `json:"type,omitempty" bson:"type,omitempty"`             // Loại dịch vụ
	AppName       constants.AppName     `json:"appname" bson:"appname"`                           // appName: facebook, instagram, tiktok,..
	FanpageID     string                `json:"fanpage_id,omitempty" bson:"fanpage_id,omitempty"` // fanpage_id, uid
	PhotoID       string                `json:"photo_id,omitempty" bson:"photo_id,omitempty"`
	PostID        string                `json:"post_id,omitempty" bson:"post_id,omitempty"`
	URLService    string                `json:"url_service,omitempty" bson:"url_service,omitempty"`
	LinkService   string                `json:"link_service,omitempty" bson:"link_service,omitempty"`
	CommentID     primitive.ObjectID    `json:"comment_id,omitempty" bson:"comment_id,omitempty"`
	ViewTime      int                   `json:"view_time,omitempty" bson:"view_time,omitempty"`
	InstagramId   string                `json:"insta_id,omitempty" bson:"insta_id,omitempty"`
	ChannelId     string                `json:"channel_id,omitempty" bson:"channel_id,omitempty"`
	ProductSearch string                `json:"product_search,omitempty" bson:"product_search,omitempty"`
	// Trace
	StartDate   *time.Time `json:"start_date,omitempty" bson:"start_date,omitempty"` // Ngày bắt đầu, dùng cho viplike service
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate *time.Time `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	ActionAt    *time.Time `json:"action_at,omitempty" bson:"action_at,omitempty"`
}

// IsExists ..
func (m ServiceLog) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// AllServicelogReq : ..
type AllServicelogReq struct {
	UID    string
	Type   string
	Status constants.ServiceLogStatus
	Limit  int
}
