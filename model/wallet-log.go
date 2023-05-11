package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WalletLog : Lưu lịch sử thay đổi ví người dùng
type WalletLog struct {
	ID          primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	Type        constants.WalletLogType `json:"type" bson:"type"`
	Token       string                  `json:"token" bson:"token"`
	Value       int                     `json:"value" bson:"value"`     // + , - phụ thuộc vào type
	Balance     int                     `json:"balance" bson:"balance"` // Tài khoản người dùng
	Transaction *TransactionWallet      `json:"transaction,omitempty" bson:"transaction,omitempty"`
	Service     *ServiceWallet          `json:"service,omitempty" bson:"service,omitempty"`
	Gift        *GiftWallet             `json:"gift,omitempty" bson:"gift,omitempty"`
	Note        string                  `json:"note,omitempty" bson:"note,omitempty"`
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
func (m WalletLog) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// AllWalletLogReq : get all wallet log
type AllWalletLogReq struct {
	Types     []constants.WalletLogType
	UserToken string
	Offset    primitive.ObjectID
	Limit     int
	FromTime  int
	ToTime    int
}

// TransactionWallet ..
type TransactionWallet struct {
	Code  string  `json:"code,omitempty" bson:"code,omitempty"`
	Value int     `json:"value,omitempty" bson:"value,omitempty"`
	Bonus float64 `json:"bonus" bson:"bonus"`
}

// ServiceWallet ..
type ServiceWallet struct {
	ServiceCode   []string              `json:"service_code,omitempty" bson:"service_code,omitempty"` // mã giao dịch, tự gen
	Type          constants.ServiceType `json:"type,omitempty" bson:"type,omitempty"`
	FanpageID     string                `json:"fanpage_id,omitempty" bson:"fanpage_id,omitempty"` // fanpage_id, uid
	PostID        string                `json:"post_id,omitempty" bson:"post_id,omitempty"`
	PhotoID       string                `json:"photo_id,omitempty" bson:"photo_id,omitempty"`
	URLService    string                `json:"url_service,omitempty" bson:"url_service,omitempty"`
	CommentID     string                `json:"comment_id,omitempty" bson:"comment_id,omitempty"`
	Price         int                   `json:"price,omitempty" bson:"price,omitempty"` // Giá
	Number        int                   `json:"number" bson:"number"`                   // Tổng sổ lượt
	NumberSuccess int                   `json:"number_success" bson:"number_success"`   // Số lượt thành công
	TotalWarranty int                   `json:"total_warranty" bson:"total_warranty"`   // Tổng số bảo hành
	Warranty      int                   `json:"warranty" bson:"warranty"`               // Thực hiện bảo hành
}

// GiftWallet ..
type GiftWallet struct {
	Code  string `json:"code,omitempty" bson:"code,omitempty"`
	Value int    `json:"value,omitempty" bson:"value,omitempty"`
}
