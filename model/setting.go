package model

import (
	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Setting : Thông tin Setting
type Setting struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AndroidLink         string             `json:"android_link" bson:"android_link"`
	IOSLink             string             `json:"ios_link" bson:"ios_link"`
	Monitoring          bool               `json:"monitoring" bson:"monitoring"`
	AndroidVersionCode  string             `json:"android_versioncode" bson:"android_versioncode"`
	FacebookApkVersion  string             `json:"facebook_apk_version" bson:"facebook_apk_version"`
	InstagramApkVersion string             `json:"instagram_apk_version" bson:"instagram_apk_version"`
	AndroidVerSionname  string             `json:"android_versionname" bson:"android_versionname"`
	ApkTestLink         string             `json:"apk_test_link" bson:"apk_test_link"`
	ApkTestVersion      string             `json:"apk_test_version" bson:"apk_test_version"`
	Dropboxaccesstoken  string             `json:"dropboxaccesstoken" bson:"dropboxaccesstoken"`
	ApkNumbers          []ApkNumber        `json:"apk_numbers" bson:"apk_numbers"`
}

type ApkNumber struct {
	Key    constants.AppName `json:"key" bson:"key"`
	Number int               `json:"number" bson:"number"`
}

// IsExists ..
func (m Setting) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
