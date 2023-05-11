package report

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyDevice : Chứa thông tin tính tiền autofarmer theo device, action theo ngày
type DailyPay struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token         string             `json:"token" bson:"token"` // autofarmer token
	Status        string             `json:"status" bson:"status"`
	Amount        int                `json:"amount" bson:"amount"`
	AccountName   string             `json:"account_name,omitempty" bson:"account_name,omitempty"`
	Username      string             `json:"username,omitempty" bson:"username,omitempty"`
	ResonanceCode string             `json:"resonance_code,omitempty" bson:"resonance_code,omitempty"`
	BankName      string             `json:"bank_name,omitempty" bson:"bank_name,omitempty"`
	BankNumber    string             `json:"bank_number,omitempty" bson:"bank_number,omitempty"`
	CreatedDate   *time.Time         `json:"created_date" bson:"created_date"`
	UpdatedDate   *time.Time         `json:"updated_date" bson:"updated_date"`
}

type AllDailyPayReq struct {
	Token     string
	Status    string
	StartTime *time.Time
	EndTime   *time.Time
}

// IsExists ..
func (m DailyPay) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
