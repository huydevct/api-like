package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceLog : Lưu thông tin service log
type ServiceLog struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Type        constants.ServiceType `json:"type,omitempty" bson:"type,omitempty"`                 // Loại dịch vụ
	ServiceCode string                `json:"service_code,omitempty" bson:"service_code,omitempty"` // mã giao dịch, tự gen
	UID         string                `json:"uid,omitempty" bson:"uid,omitempty"`
	Status      int                   `json:"status,omitempty" bson:"status,omitempty"`
	Kind        int                   `json:"kind,omitempty" bson:"kind,omitempty"`
	FanpageID   string                `json:"fanpage_id,omitempty" bson:"fanpage_id,omitempty"`
	IMEI        string                `json:"IMEI,omitempty" bson:"IMEI,omitempty"`
	Token       string                `json:"token,omitempty" bson:"token,omitempty"`
	Price       int                   `json:"price,omitempty" bson:"price,omitempty"`
	CreatedAt   int                   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	WeekYear    string                `json:"weekYear,omitempty" bson:"weekYear,omitempty"`
	Week        int                   `json:"week,omitempty" bson:"week,omitempty"`
	Year        int                   `json:"year,omitempty" bson:"year,omitempty"`
	PostID      string                `json:"postID,omitempty" bson:"postID,omitempty"`
	Liked       bool                  `json:"liked" bson:"liked"`
	ActionAt    *time.Time            `json:"action_at,omitempty" bson:"action_at,omitempty"`
}
