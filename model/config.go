package model

import "app/constants"

//Config : Thông tin Config
type Config struct {
	TimeOut             string                    `json:"timeout" bson:"timeout"`
	DebugMode           constants.DebugModeStatus `json:"debug_mode" bson:"debug_mode"`
	Reset3G             string                    `json:"reset3g" bson:"reset3g"`
	VpnLoading          string                    `json:"vpn_loading" bson:"vpn_loading"`
	UserType            string                    `json:"user_type" bson:"user_type"`
	Balance             int                       `json:"balance" bson:"balance"`
	Connectivity        string                    `json:"connectivity" bson:"connectivity"`
	Monitoring          bool                      `json:"monitoring" bson:"monitoring"`
	AndroidVersionCode  string                    `json:"android_versioncode" bson:"android_versioncode"`
	FacebookApkVersion  string                    `json:"facebook_apk_version" bson:"facebook_apk_version"`
	FacebookApkNumber   int                       `json:"facebook_apk_number" bson:"facebook_apk_number"`
	InstagramApkVersion string                    `json:"instagram_apk_version" bson:"instagram_apk_version"`
	InstagramApkNumber  int                       `json:"instagram_apk_number" bson:"instagram_apk_number"`
	TiktokApkNumber     int                       `json:"tiktok_apk_number" bson:"tiktok_apk_number"`
	ShareLiveStream     bool                      `json:"share_live_stream" bson:"share_live_stream"`
	User100App          bool                      `json:"user_100_app" bson:"user_100_app"`
	IsReg               bool                      `json:"is_reg" bson:"is_reg"`
	AndroidVerSionname  string                    `json:"android_versionname" bson:"android_versionname"`
	ApkTestLink         string                    `json:"apk_test_link" bson:"apk_test_link"`
	ApkTestVersion      string                    `json:"apk_test_version" bson:"apk_test_version"`
	Dropboxaccesstoken  string                    `json:"dropboxaccesstoken" bson:"dropboxaccesstoken"`
	// Thông tin device
	DeviceInfo *DeviceInfo `json:"device_info,omitempty" bson:"device_info,omitempty"`
}
