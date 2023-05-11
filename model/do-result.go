package model

import (
	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ViplikeDataItem : ..
type ViplikeDataItem struct {
	PostID string `json:"postID"`
	Liked  bool   `json:"liked"`
}

// DoResultPublish : data push to rabbtit
type DoResultPublish struct {
	CloneInfo      CloneInfo                `json:"clone_info"`
	DeviceID       string                   `json:"device_id"`
	MacAddress     string                   `json:"mac_address"`
	ServiceCode    string                   `json:"service_code"`
	Action         string                   `json:"action"`
	ViplikeData    []ViplikeDataItem        `json:"viplike_date"`
	SharePostID    primitive.ObjectID       `json:"share_post_id"`
	ShareCommentID primitive.ObjectID       `json:"share_comment_id"`
	ServiceLog     ServiceLog               `json:"service_log"`
	Status         constants.DoResultStatus `json:"status"`
	Reason         string                   `json:"reason"`
}
