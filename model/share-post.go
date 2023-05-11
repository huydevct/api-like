package model

import (
	"time"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllCommentReq : get all comment
type AllSharePostReq struct {
	Status []constants.CommentStatus
	Name   string
	Token  string
	Offset primitive.ObjectID
	Limit  int
}

// SharePost : Lưu thông tin các bài share
type SharePost struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token       string              `json:"token" bson:"token"`
	ServiceCode string              `json:"service_code" bson:"service_code"`
	LinkService string              `json:"link_service" bson:"link_service"`
	NumberShare int                 `json:"number_share" bson:"number_share"`
	//
	CreatedUser *primitive.ObjectID `json:"created_user,omitempty" bson:"created_user,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// SharePost : Lưu log các clone đã share
type SharePostLog struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token       string             `json:"token" bson:"token"`
	ServiceCode string             `json:"service_code" bson:"service_code"`
	LinkService string             `json:"link_service" bson:"link_service"`
	Comment     string             `json:"comment" bson:"comment"`
	CloneInfo   CloneInfo          `json:"clone_info" bson:"clone_info"`
	//
	CreatedUser *primitive.ObjectID `json:"created_user,omitempty" bson:"created_user,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// IsExists ..
func (m SharePost) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
