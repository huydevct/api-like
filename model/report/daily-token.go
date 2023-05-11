package report

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyToken : Chứa thông tin tính tiền autofarm theo action theo ngày
type DailyToken struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token                  string             `json:"token" bson:"token"`                                       // autofarmer token
	TotalLikePage          int                `json:"total_likepage" bson:"total_likepage"`                     // Tổng số action likepage
	TotalFollow            int                `json:"total_follow" bson:"total_follow"`                         // Tổng số action follow
	TotalBuffLike          int                `json:"total_bufflike" bson:"total_bufflike"`                     // Tổng số action bufflike
	TotalBuffComment       int                `json:"total_buffcomment" bson:"total_buffcomment"`               // Tổng số action buffcomment
	TotalVipLike           int                `json:"total_viplike" bson:"total_viplike"`                       // Tổng số action viplike
	TotalVipLikeNew        int                `json:"total_viplike_new" bson:"total_viplike_new"`               // Tổng số action viplikenew
	TotalVipComment        int                `json:"total_vipcomment" bson:"total_vipcomment"`                 // Tổng số action vipcomment
	TotalFollowLow         int                `json:"total_follow_low" bson:"total_follow_low"`                 // Tổng số action follow low
	TotalFollowNormal      int                `json:"total_follow_normal" bson:"total_follow_normal"`           // Tổng số action follow normal
	TotalFollowHigh        int                `json:"total_follow_high" bson:"total_follow_high"`               // Tổng số action follow high
	TotalFollsowFast       int                `json:"total_follow_fast" bson:"total_follow_fast"`               // Tổng số action follow fast
	TotalLikePageLow       int                `json:"total_likepage_low" bson:"total_likepage_low"`             // Tổng số action likepage low
	TotalLikePageNormal    int                `json:"total_likepage_normal" bson:"total_likepage_normal"`       // Tổng số action likepage normal
	TotalLikePageHigh      int                `json:"total_likepage_high" bson:"total_likepage_high"`           // Tổng số action likepage high
	TotalLikePageFast      int                `json:"total_likepage_fast" bson:"total_likepage_fast"`           // Tổng số action likepage fast
	TotalFollowWarranty    int                `json:"total_follow_warranty" bson:"total_follow_warranty"`       // Tổng số action follow warranty
	TotalLikepageWarranty  int                `json:"total_likepage_warranty" bson:"total_likepage_warranty"`   // Tổng số action vipcomment
	LikePageAmount         int                `json:"likepage_amount" bson:"likepage_amount"`                   // Tổng tiền action likepage
	FollowAmount           int                `json:"follow_amount" bson:"follow_amount"`                       // Tổng tiền action follow
	BuffLikeAmount         int                `json:"bufflike_amount" bson:"bufflike_amount"`                   // Tổng tiền action bufflike
	BuffCommentAmount      int                `json:"buffcomment_amount" bson:"buffcomment_amount"`             // Tổng tiền action buffcomment
	VipLikeAmount          int                `json:"viplike_amount" bson:"viplike_amount"`                     // Tổng tiền action viplike
	VipLikeNewAmount       int                `json:"viplike_new_amount" bson:"viplike_new_amount"`             // Tổng tiền action viplikenew
	VipCommentAmount       int                `json:"vipcomment_amount" bson:"vipcomment_amount"`               // Tổng tiền action vipcomment
	FollowWarrantyAmount   int                `json:"follow_warranty_amount" bson:"follow_warranty_amount"`     // Tổng tiền action follow warranty
	LikepageWarrantyAmount int                `json:"likepage_warranty_amount" bson:"likepage_warranty_amount"` // Tổng tiền action vipcomment
	FollowLowAmount        int                `json:"follow_low_amount" bson:"follow_low_amount"`               // Tổng tiền action follow low
	FollowNormalAmount     int                `json:"follow_normal_amount" bson:"follow_normal_amount"`         // Tổng tiền action follow normal
	FollowHighAmount       int                `json:"follow_high_amount" bson:"follow_high_amount"`             // Tổng tiền action follow high
	FollowFastAmount       int                `json:"follow_fast_amount" bson:"follow_fast_amount"`             // Tổng tiền action follow fast
	LikePageLowAmount      int                `json:"likepage_low_amount" bson:"likepage_low_amount"`           // Tổng tiền action likepage low
	LikePageNormalAmount   int                `json:"likepage_normal_amount" bson:"likepage_normal_amount"`     // Tổng tiền action likepage normal
	LikePageHighAmount     int                `json:"likepage_high_amount" bson:"likepage_high_amount"`         // Tổng tiền action likepage high
	LikePageFastAmount     int                `json:"likepage_fast_amount" bson:"likepage_fast_amount"`         // Tổng tiền action likepage fast
	CreatedDate            *time.Time         `json:"created_date" bson:"created_date"`
	UpdatedDate            *time.Time         `json:"updated_date" bson:"updated_date"`
}

type AllDailyTokenReq struct {
	Token     string
	StartTime *time.Time
	EndTime   *time.Time
}
