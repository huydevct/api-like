package report

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyFee : Chứa thông tin insert vào bảng report_daily_fee của user autofarmer
type DailyFee struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token    string             `json:"token" bson:"token"` // autofarmer token
	DeviceID string             `json:"device_id" bson:"device_id"`
	// Facebook Action
	// follow
	TotalFacebookFollow         int `json:"total_facebook_follow" bson:"total_facebook_follow"`
	TotalFacebookFollowNormal   int `json:"total_facebook_follow_normal" bson:"total_facebook_follow_normal"`
	TotalFacebookFollowLow      int `json:"total_facebook_follow_low" bson:"total_facebook_follow_low"`
	TotalFacebookFollowHigh     int `json:"total_facebook_follow_high" bson:"total_facebook_follow_high"`
	TotalFacebookFollowFast     int `json:"total_facebook_follow_fast" bson:"total_facebook_follow_fast"`
	TotalFacebookFollowWarranty int `json:"total_facebook_follow_warranty" bson:"total_facebook_follow_warranty"`
	// likepage
	TotalFacebookLikepage         int `json:"total_facebook_likepage" bson:"total_facebook_likepage"`
	TotalFacebookLikePageNormal   int `json:"total_facebook_likepage_normal" bson:"total_facebook_likepage_normal"`
	TotalFacebookLikePageLow      int `json:"total_facebook_likepage_low" bson:"total_facebook_likepage_low"`
	TotalFacebookLikePageHigh     int `json:"total_facebook_likepage_high" bson:"total_facebook_likepage_high"`
	TotalFacebookLikePageFast     int `json:"total_facebook_likepag_fast" bson:"total_facebook_likepag_fast"`
	TotalFacebookLikePageWarranty int `json:"total_facebook_likepage_warranty" bson:"total_facebook_likepage_warranty"`
	//
	TotalFacebookViplikeService int `json:"total_facebook_viplikeService" bson:"total_facebook_viplikeService"`
	TotalFacebookBufflike       int `json:"total_facebook_bufflike" bson:"total_facebook_bufflike"`
	TotalFacebookBuffcomment    int `json:"total_facebook_buffcomment" bson:"total_facebook_buffcomment"`
	TotalFacebookVipcomment     int `json:"total_facebook_vipcomment" bson:"total_facebook_vipcomment"`

	// Autofarmer Action
	TotalAutofarmerFeed             int `json:"total_autofarmer_Feed" bson:"total_autofarmer_Feed"`
	TotalAutofarmerFeedLike         int `json:"total_autofarmer_FeedLike" bson:"total_autofarmer_FeedLike"`
	TotalAutofarmerWatch            int `json:"total_autofarmer_Watch" bson:"total_autofarmer_Watch"`
	TotalAutofarmerWatchLike        int `json:"total_autofarmer_WatchLike" bson:"total_autofarmer_WatchLike"`
	TotalAutofarmerPageFeed         int `json:"total_autofarmer_PageFeed" bson:"total_autofarmer_PageFeed"`
	TotalAutofarmerGroupFeed        int `json:"total_autofarmer_GroupFeed" bson:"total_autofarmer_GroupFeed"`
	TotalAutofarmerGroupLike        int `json:"total_autofarmer_GroupLike" bson:"total_autofarmer_GroupLike"`
	TotalAutofarmerJoinGroup        int `json:"total_autofarmer_JoinGroup" bson:"total_autofarmer_JoinGroup"`
	TotalAutofarmerFriendFeed       int `json:"total_autofarmer_FriendFeed" bson:"total_autofarmer_FriendFeed"`
	TotalAutofarmerFriendLike       int `json:"total_autofarmer_FriendLike" bson:"total_autofarmer_FriendLike"`
	TotalAutofarmerAddFriendSuggest int `json:"total_autofarmer_AddFriendSuggest" bson:"total_autofarmer_AddFriendSuggest"`
	TotalAutofarmerAddFriendUID     int `json:"total_autofarmer_AddFriendUID" bson:"total_autofarmer_AddFriendUID"`
	TotalAutofarmerConfirmFriend    int `json:"total_autofarmer_ConfirmFriend" bson:"total_autofarmer_ConfirmFriend"`
	TotalAutofarmerShareNow         int `json:"total_autofarmer_ShareNow" bson:"total_autofarmer_ShareNow"`
	TotalAutofarmerShareGroup       int `json:"total_autofarmer_ShareGroup" bson:"total_autofarmer_ShareGroup"`
	TotalAutofarmerPostStatusImage  int `json:"total_autofarmer_PostStatusImage" bson:"total_autofarmer_PostStatusImage"`
	TotalAutofarmerFollow           int `json:"total_autofarmer_Follow" bson:"total_autofarmer_Follow"`
	TotalAutofarmerCommentTag       int `json:"total_autofarmer_CommentTag" bson:"total_autofarmer_CommentTag"`
	TotalAutofarmerSearch           int `json:"total_autofarmer_Search" bson:"total_autofarmer_Search"`
	TotalAutofarmerFeedView         int `json:"total_autofarmer_FeedView" bson:"total_autofarmer_FeedView"`
	TotalAutofarmerSubcrible        int `json:"total_autofarmer_Subcrible" bson:"total_autofarmer_Subcrible"`
	TotalAutofarmerComment          int `json:"total_autofarmer_Comment" bson:"total_autofarmer_Comment"`
	TotalAutofarmerLike             int `json:"total_autofarmer_Like" bson:"total_autofarmer_Like"`
	TotalAutofarmerDisLike          int `json:"total_autofarmer_DisLike" bson:"total_autofarmer_DisLike"`
	// Youtube Action
	TotalYoutubeSubscribe int `json:"total_youtube_Subscribe" bson:"total_youtube_Subscribe"`
	TotalYoutubeComment   int `json:"total_youtube_Comment" bson:"total_youtube_Comment"`
	TotalYoutubeLike      int `json:"total_youtube_Like" bson:"total_youtube_Like"`
	TotalYoutubeDisLike   int `json:"total_youtube_DisLike" bson:"total_youtube_DisLike"`
	TotalYoutubeSearch    int `json:"total_youtube_Search" bson:"total_youtube_Search"`
	// Instagram Action
	TotalInstagramFollow     int        `json:"total_instagram_Follow" bson:"total_instagram_Follow"`
	TotalInstagramUnFollow   int        `json:"total_instagram_UnFollow" bson:"total_instagram_UnFollow"`
	TotalInstagramBuffLike   int        `json:"total_instagram_BuffLike" bson:"total_instagram_BuffLike"`
	TotalInstagramComment    int        `json:"total_instagram_Comment" bson:"total_instagram_Comment"`
	TotalInstagramCommentTag int        `json:"total_instagram_CommentTag" bson:"total_instagram_CommentTag"`
	TotalInstagramSearch     int        `json:"total_instagram_Search" bson:"total_instagram_Search"`
	Total                    int        `json:"total" bson:"total"`             // Tổng số loại action
	TotalPrice               int        `json:"total_price" bson:"total_price"` // Tổng số tiền của mỗi device
	CreatedDate              *time.Time `json:"created_date" bson:"created_date"`
	UpdatedDate              *time.Time `json:"updated_date" bson:"updated_date"`
}

// AllDailyFeeReq : Chứa thông tin request lấy all data trong bảng report_daily_fee
type AllDailyFeeReq struct {
	Token     string
	StartTime *time.Time
	EndTime   *time.Time
	Limit     int
	Page      int
}
