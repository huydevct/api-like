package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service : chứa thông tin service
type Service struct {
	ID                 primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	ServiceCode        string                  `json:"service_code,omitempty" bson:"service_code,omitempty"`                 // mã giao dịch, tự gen
	TransactionCode    string                  `json:"transaction_code,omitempty" bson:"transaction_code,omitempty"`         // mã thanh toán, dùng cho đại lý
	ViplikeServiceCode string                  `json:"viplike_service_code,omitempty" bson:"viplike_service_code,omitempty"` // Mã dịch vụ dùng cho viplikeService
	Type               constants.ServiceType   `json:"type,omitempty" bson:"type,omitempty"`                                 // Loại dịch vụ
	AppName            constants.AppName       `json:"appname" bson:"appname"`                                               // appName: facebook, instagram, tiktok,..
	Kind               constants.ServiceKind   `json:"kind,omitempty" bson:"kind,omitempty"`                                 // Gói thuờng, hay gói bảo hành
	Status             constants.ServiceStatus `json:"status,omitempty" bson:"status,omitempty"`
	ReasonReport       string                  `json:"reason_report,omitempty" bson:"reason_report,omitempty"`
	Token              string                  `json:"token,omitempty" bson:"token,omitempty"`
	UserID             *primitive.ObjectID     `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Username           string                  `json:"username,omitempty" bson:"username,omitempty"`
	Fullname           string                  `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Note               string                  `json:"note,omitempty" bson:"note,omitempty"`

	UserToken      string `json:"user_token,omitempty" bson:"user_token,omitempty"`
	ServiceCodeOld string `json:"service_code_old,omitempty" bson:"service_code_old,omitempty"`
	LikeStart      int    `json:"likes_start,omitempty" bson:"likes_start,omitempty"`
	LikeEnd        int    `json:"likes_end,omitempty" bson:"likes_end,omitempty"`
	FollowStart    int    `json:"follows_start,omitempty" bson:"follows_start,omitempty"`
	FollowEnd      int    `json:"follows_end,omitempty" bson:"follows_end,omitempty"`
	// Follow, LikePage, VipLike
	FanpageID string `json:"fanpage_id,omitempty" bson:"fanpage_id,omitempty"` // fanpage_id, uid
	// BuffLike
	PostID      string `json:"post_id,omitempty" bson:"post_id,omitempty"`
	PhotoID     string `json:"photo_id,omitempty" bson:"photo_id,omitempty"`
	URLService  string `json:"url_service,omitempty" bson:"url_service,omitempty"`
	LinkService string `json:"link_service,omitempty" bson:"link_service,omitempty"`
	// BuffComment, VipComment
	CommentID primitive.ObjectID `json:"comment_id,omitempty" bson:"comment_id,omitempty"`
	// BuffView
	ViewTime int `json:"view_time,omitempty" bson:"view_time,omitempty"`
	// Instagram follow
	InstagramId string `json:"insta_id,omitempty" bson:"insta_id,omitempty"`
	// Youtube
	ChannelId string `json:"channel_id,omitempty" bson:"channel_id,omitempty"`
	// Shopee, Lazada
	ProductSearch string `json:"product_search,omitempty" bson:"product_search,omitempty"`
	// Số lượng
	Price            int `json:"price,omitempty" bson:"price,omitempty"` // Giá
	Number           int `json:"number" bson:"number"`                   // Tổng sổ lượt
	NumberSuccess    int `json:"number_success" bson:"number_success"`
	PriceInt         int `json:"price_int" bson:"price_int"`                   // Giá int
	NumberRest       int `json:"NumberRest" bson:"NumberRest"`                 // gia tri con lai
	NumberInt        int `json:"number_int" bson:"number_int"`                 // Tổng sổ lượt int
	NumberSuccessInt int `json:"number_success_int" bson:"number_success_int"` // Số lượt thành công int
	TotalWarranty    int `json:"total_warranty" bson:"total_warranty"`         // Tổng số bảo hành
	Warranty         int `json:"warranty" bson:"warranty"`                     // Thực hiện bảo hành
	CreatedAt        int `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateTime       int `json:"updateTime,omitempty" bson:"updateTime,omitempty"`
	TimeSuccess      int `json:"TimeSuccess,omitempty" bson:"TimeSuccess,omitempty"`
	// Report
	NumberReport int `json:"number_report" bson:"number_report"` // tổng số number report
	// NumberDeff    int    `json:"number_deff" bson:"number_deff"`                 // Số lượt còn lại
	// Extra data
	Data        []ViplikeItemOld `json:"data" bson:"data"`
	DataViplike ViplikeItem      `json:"data_viplike" bson:"data_viplike"`
	AliveDays   int              `json:"alive_days,omitempty" bson:"alive_days,omitempty"` // Số ngày chạy gói viplike
	DayLeft     int              `json:"day_left" bson:"day_left"`                         // Ngày còn lại của gói viplike
	StartDate   *time.Time       `json:"start_date,omitempty" bson:"start_date,omitempty"` // Ngày bắt đầu, dùng cho viplike service
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
func (m Service) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// ViplikeItem : Lưu số lượt like của gói viplike, số lượng post có thể thay đổi
type ViplikeItem struct {
	Post1 int `json:"p1" bson:"p1"`
	Post2 int `json:"p2" bson:"p2"`
	Post3 int `json:"p3" bson:"p3"`
	Post4 int `json:"p4" bson:"p4"`
	Post5 int `json:"p5" bson:"p5"`
}

// ViplikeItemOld : Lưu số lượt like của gói viplike old, số lượng post có thể thay đổi
type ViplikeItemOld struct {
	PostID string `json:"postID" bson:"postID"`
	Count  int    `json:"count" bson:"count"`
}

// CreateServicePublishReq : Data push to rabbit khi tạo gói
type CreateServicePublishReq struct {
	Services []Service `json:"services"`
}
