package model

import (
	"time"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Package : Lưu thông tin gói nạp tiền
type Package struct {
	ID              primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Money           int                    `json:"money" bson:"money"` // Giá trị gói
	Bonus           float64                `json:"bonus" bson:"bonus"` // Thưởng
	Type            constants.PackageType  `json:"type" bson:"type"`   // Loai gói
	Status          constants.CommonStatus `json:"status" bson:"status"`
	IsLikeSub       bool                   `json:"is_like_sub" bson:"is_like_sub"`
	IsInstagram     bool                   `json:"is_instagram" bson:"is_instagram"`
	IsYoutube       bool                   `json:"is_youtube" bson:"is_youtube"`
	IsReg           bool                   `json:"is_reg" bson:"is_reg"`
	ShareLiveStream bool                   `json:"share_live_stream" bson:"share_live_stream"`
	User100App      bool                   `json:"user_100_app" bson:"user_100_app"`
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
func (m Package) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// AllPackageReq : get all package
type AllPackageReq struct {
	Status []constants.CommonStatus
	Offset primitive.ObjectID
	Limit  int
}
