package model

import (
	"time"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllCommentReq : get all comment
type AllShareCommentReq struct {
	Status []constants.CommentStatus
	Name   string
	Token  string
	Offset primitive.ObjectID
	Limit  int
}

// SharePost : Lưu thông tin các bài share
type ShareComment struct {
	ID            *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token         string              `json:"token" bson:"token"`
	ServiceCode   string              `json:"service_code" bson:"service_code"`
	Comment       string              `json:"comment" bson:"comment"`
	NumberComment int                 `json:"number_comment" bson:"number_comment"`
	//
	CreatedUser *primitive.ObjectID `json:"created_user,omitempty" bson:"created_user,omitempty"`
	CreatedDate *time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedDate *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// IsExists ..
func (m ShareComment) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
