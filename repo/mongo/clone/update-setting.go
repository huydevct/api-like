package clone

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateCloneSettingReq Cập nhật thông tin cài đặt của clone
type UpdateCloneSettingReq struct {
	CloneID          primitive.ObjectID
	Token            string
	SettingSecretkey *bool // true: clone đã cài secretkey
	SettingAvatar    *bool // true: clone đã cài avatar
	SettingCover     *bool // true: clone đã cài background cover
	SettingLang      *bool // true: clone đã cài lang
}

// UpdateSetting : Cập nhật thông tin setting của clone
// 1. condition: id, token
func (r Repo) UpdateSetting(req UpdateCloneSettingReq) (err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	// set condition
	condition := bson.M{}
	condition["_id"] = req.CloneID
	condition["token"] = req.Token

	// set update
	now := time.Now()
	update := bson.M{
		"updated_date": now,
	}
	if req.SettingSecretkey != nil {
		update["setting_secretkey"] = *req.SettingSecretkey
	}
	if req.SettingAvatar != nil {
		update["setting_avatar"] = *req.SettingAvatar
	}
	if req.SettingCover != nil {
		update["setting_cover"] = *req.SettingCover
	}
	if req.SettingLang != nil {
		update["setting_lang"] = *req.SettingLang
	}
	updates := bson.M{
		"$set": update,
	}

	updateResult, err := collection.UpdateOne(ctx,
		condition,
		updates)
	if err != nil {
		return
	}
	if updateResult.MatchedCount == 0 {
		err = fmt.Errorf("Update fail, not map condition")
	}

	return
}
