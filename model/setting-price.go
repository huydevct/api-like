package model

import (
	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SettingPrice : setting giá cho autofarmer
type SettingPrice struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Appname  string             `json:"appname" bson:"appname"`
	Settings []SettingPriceItem `json:"settings" bson:"settings"`
}

// SettingPriceItem : giá cho từng key
type SettingPriceItem struct {
	Key              string             `json:"key" bson:"key"`
	Type             constants.AppName  `json:"type" bson:"type"`
	Price            int                `json:"price" bson:"price"`
	Limit            int                `json:"limit" bson:"limit"`
	UpdateNumberRest []UpdateNumberRest `json:"update_number_rest" bson:"update_number_rest"`
}

// UpdateNumberRest:  cấu hình number rest
type UpdateNumberRest struct {
	Number     int `json:"number" bson:"number"`
	NumberRest int `json:"number_rest" bson:"number_rest"`
}

// ByNumber
type ByNumber []UpdateNumberRest

func (a ByNumber) Len() int           { return len(a) }
func (a ByNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNumber) Less(i, j int) bool { return a[i].Number < a[j].Number }

// IsExists ..
func (m *SettingPrice) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
