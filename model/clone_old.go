package model

import (
	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CloneOldInfo : Chứa thông tin old clone
type CloneOldInfo struct {
	ID               primitive.ObjectID     `json:"_id"`
	Email            string                 `json:"email"` // (unit)
	UID              string                 `json:"uid"`   // (unit)
	Token            string                 `json:"token"`
	IMEI             string                 `json:"IMEI"` // DeviceID
	AliveStatus      constants.AliveStatus  `json:"alive_status"`
	Secretkey        string                 `json:"fa"`
	Password         string                 `json:"password"`
	Passmail         string                 `json:"passmail"`
	PhoneNumber      string                 `json:"phone_number"`
	Language         string                 `json:"language"`
	Country          constants.CloneCountry `json:"country"`
	AppName          constants.AppName      `json:"appname"`
	IsAutofarmer     interface{}            `json:"is_autofarmer"`
	SettingSecretkey interface{}            `json:"Setting2FA"`
	SettingAvatar    interface{}            `json:"avatar"`
	SettingCover     interface{}            `json:"cover"`
	SettingLang      interface{}            `json:"SettingLang"`
	MobileEncode     bool                   `json:"mz"`
	CGIEncode        bool                   `json:"cz"`
}

// IsExists ..
func (m CloneOldInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
