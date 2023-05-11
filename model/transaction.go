package model

import (
	"time"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction : Lưu thông tin nạp tiền
type Transaction struct {
	ID              primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Code            string                `json:"code,omitempty" bson:"code,omitempty"`         // Mã tự gen 11 ký tự
	UserID          *primitive.ObjectID   `json:"user_id,omitempty" bson:"user_id,omitempty"`   // Sdt
	Username        string                `json:"username,omitempty" bson:"username,omitempty"` // Sdt
	Fullname        string                `json:"fullname,omitempty" bson:"fullname,omitempty"` // Tên nhân viên
	Token           string                `json:"token,omitempty" bson:"token,omitempty"`       // User token
	Status          constants.TransStatus `json:"status,omitempty" bson:"status,omitempty"`     // Pending | Active
	Value           interface{}           `json:"value,omitempty" bson:"value,omitempty"`
	IsLikeSub       bool                  `json:"is_like_sub,omitempty" bson:"is_like_sub,omitempty"`
	IsInstagram     bool                  `json:"is_instagram,omitempty" bson:"is_instagram,omitempty"`
	IsYoutube       bool                  `json:"is_youtube,omitempty" bson:"is_youtube,omitempty"`
	IsReg           bool                  `json:"is_reg,omitempty" bson:"is_reg,omitempty"`
	ShareLiveStream bool                  `json:"share_live_stream,omitempty" bson:"share_live_stream,omitempty"`
	User100App      bool                  `json:"user_100_app,omitempty" bson:"user_100_app,omitempty"`
	ValueInt        int                   `json:"value_int,omitempty" bson:"value_int,omitempty"`
	Bonus           float64               `json:"bonus" bson:"bonus"`
	CreatedAt       int                   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	//
	CreatedIP   string     `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
}

// IsExists ..
func (m Transaction) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// AllTransactionReq : get all transaction
type AllTransactionReq struct {
	Status    []constants.TransStatus
	Code      string
	UserToken string
	Username  string
	Offset    primitive.ObjectID
	Limit     int
}
