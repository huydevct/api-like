package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActiveService : Lưu danh sách các gói service đang active
type ActiveService struct {
	ID                 primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	ServiceCode        string                  `json:"service_code" bson:"service_code"`
	ViplikeServiceCode string                  `json:"viplike_service_code,omitempty" bson:"viplike_service_code,omitempty"` // Mã dịch vụ dùng cho viplikeService
	Status             constants.ServiceStatus `json:"status" bson:"status"`
	Type               string                  `json:"type,omitempty" bson:"type,omitempty"`
	// info
	FanpageID   string             `json:"fanpage_id,omitempty" bson:"fanpage_id,omitempty"` // fanpage_id, uid
	URLService  string             `json:"url_service,omitempty" bson:"url_service,omitempty"`
	LinkService string             `json:"link_service,omitempty" bson:"link_service,omitempty"`
	CommentID   primitive.ObjectID `json:"comment_id,omitempty" bson:"comment_id,omitempty"`
	// Số lượng
	Number        int `json:"number" bson:"number"` // Tổng sổ lượt
	NumberSuccess int `json:"number_success" bson:"number_success"`
	// trace
	StartDate     *time.Time `json:"start_date,omitempty" bson:"start_date,omitempty"` // ngày bắt đầu chaỵ
	CreatedDate   *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate   *time.Time `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	GetRandomTime *time.Time `json:"get_random_time,omitempty" bson:"get_random_time,omitempty"`
}

// IsExists ..
func (m ActiveService) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
