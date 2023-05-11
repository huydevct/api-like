package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllCommentContentReq : get all comment content
type AllCommentContentReq struct {
	CommentID string
	Token     string
	Offset    primitive.ObjectID
	Limit     int
}

// CommentContent : Lưu nội dung bình luận
type CommentContent struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CommentID string             `json:"comment_id" bson:"comment_id"`
	Content   string             `json:"content" bson:"content"`
	Token     string             `json:"token" bson:"token"`
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
func (m CommentContent) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
