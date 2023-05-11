package model

import (
	"app/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AllActionProfileReq : get all action
type AllActionProfileReq struct {
	Name     string
	Offset   primitive.ObjectID
	AppName  string
	Token    string
	Template constants.TemplateAction
	Status   []constants.CommonStatus
	Limit    int
}

// ActionProfile : thông tin action profile
type ActionProfile struct {
	ID            primitive.ObjectID       `json:"id" bson:"_id,omitempty"`
	Type          string                   `json:"type"  bson:"type"`
	Name          string                   `json:"name" bson:"name"`
	AppName       string                   `json:"appname" bson:"appname"`
	Token         string                   `json:"token" bson:"token"`
	Template      constants.TemplateAction `json:"template" bson:"template"`
	Status        constants.CommonStatus   `json:"status" bson:"status"`
	Actions       [][]Action               `json:"actions" bson:"actions,omitempty"`
	ActionDefault []Action                 `json:"action_default,omitempty" bson:"action_default,omitempty"`
	CreatedDate   *time.Time               `json:"created_date,omitempty" bson:"created_date,omitempty"`
}

//Action : thông tin action
type Action struct {
	Action   constants.AutofarmerAction `json:"action" bson:"action"`
	ReRun    int                        `json:"re_run" bson:"re_run"`
	Quantity int                        `json:"quantity" bson:"quantity"`
	IsRandom bool                       `json:"is_random" bson:"is_random"`
}

// IsExists ..
func (m ActionProfile) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
