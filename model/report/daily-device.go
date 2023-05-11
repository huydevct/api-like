package report

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyDevice : Chứa thông tin tính tiền autofarmer theo device, action theo ngày
type DailyDevice struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token    string             `json:"token" bson:"token"` // autofarmer token
	DeviceID string             `json:"device_id" bson:"device_id"`
	// likepage
	TotalLikePage         int `json:"total_likepage" bson:"total_likepage"`                   // Tổng số action likepage
	TotalLikePageNormal   int `json:"total_likepage_normal" bson:"total_likepage_normal"`     // Tổng số action follow normal
	TotalLikePageLow      int `json:"total_likepage_low" bson:"total_likepage_low"`           // Tổng số action likepage low
	TotalLikePageHigh     int `json:"total_likepage_high" bson:"total_likepage_high"`         // Tổng số action likepage high
	TotalLikePageFast     int `json:"total_likepage_fast" bson:"total_likepage_fast"`         // Tổng số action likepage fast
	TotalLikepageWarranty int `json:"total_likepage_warranty" bson:"total_likepage_warranty"` // Tổng số action vipcomment

	// follow
	TotalFollow         int `json:"total_follow" bson:"total_follow"`                   // Tổng số action follow
	TotalFollowNormal   int `json:"total_follow_normal" bson:"total_follow_normal"`     // Tổng số action follow normal
	TotalFollowLow      int `json:"total_follow_low" bson:"total_follow_low"`           // Tổng số action follow low
	TotalFollowHigh     int `json:"total_follow_high" bson:"total_follow_high"`         // Tổng số action follow high
	TotalFollowFast     int `json:"total_follow_fast" bson:"total_follow_fast"`         // Tổng số action follow fast
	TotalFollowWarranty int `json:"total_follow_warranty" bson:"total_follow_warranty"` // Tổng số action follow warranty
	//
	TotalBuffLike    int        `json:"total_bufflike" bson:"total_bufflike"`       // Tổng số action bufflike
	TotalBuffComment int        `json:"total_buffcomment" bson:"total_buffcomment"` // Tổng số action buffcomment
	TotalVipLike     int        `json:"total_viplike" bson:"total_viplike"`         // Tổng số action viplike
	TotalVipLikeNew  int        `json:"total_viplike_new" bson:"total_viplike_new"` // Tổng số action viplikenew
	TotalVipComment  int        `json:"total_vipcomment" bson:"total_vipcomment"`   // Tổng số action vipcomment
	Total            int        `json:"total" bson:"total"`                         // Tổng số loại action
	TotalPrice       int        `json:"total_price" bson:"total_price"`             // Tổng số tiền của mỗi device
	CreatedDate      *time.Time `json:"created_date" bson:"created_date"`
	UpdatedDate      *time.Time `json:"updated_date" bson:"updated_date"`
}

type AllDailyDeviceReq struct {
	DeviceID  string
	Token     string
	StartTime *time.Time
	EndTime   *time.Time
	Limit     int
	Page      int
}
