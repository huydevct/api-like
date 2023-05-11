package constants

type (
	// ServiceType : lọai dịch vụ
	ServiceType string

	// SortServiceType : define loại sort
	SortServiceType string

	// AppName : appname
	AppName string

	// CloneCountry : loại clone
	CloneCountry string
	// Status : trạng thái thanh toán
	Status string
)

//  Status
const (
	StatusUnpaid   Status = "unpaid"
	StatusFinished Status = "finished"
)

// Service type
const (
	ServiceFollow           ServiceType = "follow"
	ServiceFollowLow        ServiceType = "follow_low"
	ServiceFollowHigh       ServiceType = "follow_high"
	ServiceFollowFast       ServiceType = "follow_fast"
	ServiceLikePage         ServiceType = "likepage"
	ServiceLikePageLow      ServiceType = "likepage_low"
	ServiceLikePageHigh     ServiceType = "likepage_high"
	ServiceLikePageFast     ServiceType = "likepage_fast"
	ServiceBuffReaction     ServiceType = "buff_reaction"
	ServiceBuffLike         ServiceType = "bufflike"
	ServiceBuffComment      ServiceType = "buffcomment"
	ServiceVipLike          ServiceType = "viplikeService"
	ServiceVipLikeNew       ServiceType = "viplikeServiceNew"
	ServiceVipComment       ServiceType = "vipcomment"
	ServiceFollowWarranty   ServiceType = "follow_warranty"
	ServiceLikePageWarranty ServiceType = "likepage_warranty"
	ServiceShareNow         ServiceType = "ShareNow"
	// youtube service
	ServiceYoutubeBuffLike        ServiceType = "youtube_bufflike"
	ServiceYoutubeBuffLikeLow     ServiceType = "youtube_bufflike_low"
	ServiceYoutubeBuffLikeHigh    ServiceType = "youtube_bufflike_high"
	ServiceYoutubeBuffView        ServiceType = "youtube_buffview"
	ServiceYoutubeBuffViewLow     ServiceType = "youtube_buffview_low"
	ServiceYoutubeBuffViewHigh    ServiceType = "youtube_buffview_high"
	ServiceYoutubeBuffSub         ServiceType = "youtube_buffsub"
	ServiceYoutubeBuffSubLow      ServiceType = "youtube_buffsub_low"
	ServiceYoutubeBuffSubHigh     ServiceType = "youtube_buffsub_high"
	ServiceYoutubeBuffComment     ServiceType = "youtube_buffcomment"
	ServiceYoutubeBuffCommentLow  ServiceType = "youtube_buffcomment_low"
	ServiceYoutubeBuffCommentHigh ServiceType = "youtube_buffcomment_high"
	// instagram service
	ServiceInstagramBuffLike        ServiceType = "instagram_bufflike"
	ServiceInstagramBuffLikeLow     ServiceType = "instagram_bufflike_low"
	ServiceInstagramBuffLikeHigh    ServiceType = "instagram_bufflike_high"
	ServiceInstagramBuffView        ServiceType = "instagram_buffview"
	ServiceInstagramBuffViewLow     ServiceType = "instagram_buffview_low"
	ServiceInstagramBuffViewHigh    ServiceType = "instagram_buffview_high"
	ServiceInstagramFollow          ServiceType = "instagram_follow"
	ServiceInstagramFollowLow       ServiceType = "instagram_follow_low"
	ServiceInstagramFollowHigh      ServiceType = "instagram_follow_high"
	ServiceInstagramBuffComment     ServiceType = "instagram_buffcomment"
	ServiceInstagramBuffCommentLow  ServiceType = "instagram_buffcomment_low"
	ServiceInstagramBuffCommentHigh ServiceType = "instagram_buffcomment_high"
	// app tiktok
	ServiceTikTokBuffView        ServiceType = "tiktok_buffview"
	ServiceTikTokBuffViewLow     ServiceType = "tiktok_buffview_low"
	ServiceTikTokBuffViewHigh    ServiceType = "tiktok_buffview_high"
	ServiceTikTokBuffLike        ServiceType = "tiktok_bufflike"
	ServiceTikTokBuffLikeLow     ServiceType = "tiktok_bufflike_low"
	ServiceTikTokBuffLikeHigh    ServiceType = "tiktok_bufflike_high"
	ServiceTikTokFollow          ServiceType = "tiktok_follow"
	ServiceTikTokFollowLow       ServiceType = "tiktok_follow_low"
	ServiceTikTokFollowHigh      ServiceType = "tiktok_follow_high"
	ServiceTikTokBuffComment     ServiceType = "tiktok_buffcomment"
	ServiceTikTokBuffCommentLow  ServiceType = "tiktok_buffcomment_low"
	ServiceTikTokBuffCommentHigh ServiceType = "tiktok_buffcomment_high"
	// app shopee
	ServiceShopeeBuffView        ServiceType = "shopee_buffview"
	ServiceShopeeBuffViewLow     ServiceType = "shopee_buffview_low"
	ServiceShopeeBuffViewHigh    ServiceType = "shopee_buffview_high"
	ServiceShopeeBuffEye         ServiceType = "shopee_buffeye"
	ServiceShopeeBuffEyeLow      ServiceType = "shopee_buffeye_low"
	ServiceShopeeBuffEyeHigh     ServiceType = "shopee_buffeye_high"
	ServiceShopeeFollow          ServiceType = "shopee_follow"
	ServiceShopeeFollowLow       ServiceType = "shopee_follow_low"
	ServiceShopeeFollowHigh      ServiceType = "shopee_follow_high"
	ServiceShopeeBuffComment     ServiceType = "shopee_buffcomment"
	ServiceShopeeBuffCommentLow  ServiceType = "shopee_buffcomment_low"
	ServiceShopeeBuffCommentHigh ServiceType = "shopee_buffcomment_high"
	ServiceShopeeBuffSearch      ServiceType = "shopee_buffsearch"
	ServiceShopeeBuffSearchLow   ServiceType = "shopee_buffsearch_low"
	ServiceShopeeBuffSearchHigh  ServiceType = "shopee_buffsearch_high"
	// app lazada
	ServiceLazadaBuffEye         ServiceType = "lazada_buffeye"
	ServiceLazadaBuffEyeLow      ServiceType = "lazada_buffeye_low"
	ServiceLazadaBuffEyeHigh     ServiceType = "lazada_buffeye_high"
	ServiceLazadaFollow          ServiceType = "lazada_follow"
	ServiceLazadaFollowLow       ServiceType = "lazada_follow_low"
	ServiceLazadaFollowHigh      ServiceType = "lazada_follow_high"
	ServiceLazadaBuffComment     ServiceType = "lazada_buffcomment"
	ServiceLazadaBuffCommentLow  ServiceType = "lazada_buffcomment_low"
	ServiceLazadaBuffCommentHigh ServiceType = "lazada_buffcomment_high"
	ServiceLazadaBuffSearch      ServiceType = "lazada_buffsearch"
	ServiceLazadaBuffSearchLow   ServiceType = "lazada_buffsearch_low"
	ServiceLazadaBuffSearchHigh  ServiceType = "lazada_buffsearch_high"
)

// SortServiceType :
const (
	SortTimeNumberSuccess SortServiceType = "TimeSuccess"
	SortTimeUpdated       SortServiceType = "updatedDateTime"
	SortCreatedAt         SortServiceType = "_id"
)

// AppName :
const (
	AppNameAutofarmer AppName = "autofarmer"
	AppNameFaceBook   AppName = "facebook"
	AppNameInstagram  AppName = "instagram"
	AppNameTiktok     AppName = "tiktok"
	AppNameYoutube    AppName = "youtube"
	AppNameShopee     AppName = "shopee"
	AppNameLazada     AppName = "lazada"
)

// Type : loại clone
const (
	CloneCountryVn   CloneCountry = "vn"
	CloneCountryIndo CloneCountry = "indo"
	CloneCountryEn   CloneCountry = "en"
	CloneCountryAu   CloneCountry = "au"
)

// String : convert to string
func (m CloneCountry) String() string {
	return string(m)
}

// String : convert to string
func (m AppName) String() string {
	return string(m)
}

// String : convert to string
func (m ServiceType) String() string {
	return string(m)
}

// String : convert to string
func (m SortServiceType) String() string {
	return string(m)
}

// String : convert to string
func (m Status) String() string {
	return string(m)
}
