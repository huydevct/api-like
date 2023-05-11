package model

import (
	"time"

	"app/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserInfo : Thông tin người dùng
type UserInfo struct {
	ID                   primitive.ObjectID            `json:"id" bson:"_id,omitempty"`
	Fullname             string                        `json:"fullname" bson:"fullname"`
	Username             string                        `json:"username" bson:"username"`
	Password             string                        `json:"-" bson:"password"`
	OldToken             string                        `json:"old_token,omitempty" bson:"old_token,omitempty"` // Token sync từ api4
	DeviceTotal          int                           `json:"device_total" bson:"device_total"`
	Token                string                        `json:"token" bson:"token"`
	AccountName          string                        `json:"account_name,omitempty" bson:"account_name,omitempty"`
	ResonanceCode        string                        `json:"resonance_code,omitempty" bson:"resonance_code,omitempty"`
	BankName             string                        `json:"bank_name,omitempty" bson:"bank_name,omitempty"`
	BankNumber           string                        `json:"bank_number,omitempty" bson:"bank_number,omitempty"`
	Balance              int                           `json:"balance" bson:"balance"`
	TimeOut              string                        `json:"timeout" bson:"timeout"`
	DebugMode            constants.DebugModeStatus     `json:"debug_mode" bson:"debug_mode"`
	ActionProfileDefault []ActionProfileDefault        `json:"action_profile_default" bson:"action_profile_default"`
	Reset3G              string                        `json:"Reset3G" bson:"Reset3G"`
	Connectivity         string                        `json:"connectivity" bson:"connectivity"`
	UserType             string                        `json:"user_type" bson:"user_type"`
	ShareLiveStream      bool                          `json:"share_live_stream" bson:"share_live_stream"`
	User100App           bool                          `json:"user_100_app" bson:"user_100_app"`
	CheckApp             bool                          `json:"check_app" bson:"check_app"`
	IsReg                bool                          `json:"is_reg" bson:"is_reg"`
	IsLikeSub            bool                          `json:"is_like_sub,omitempty" bson:"is_like_sub"`
	IsInstagram          bool                          `json:"is_instagram,omitempty" bson:"is_instagram"`
	IsYoutube            bool                          `json:"is_youtube,omitempty" bson:"is_youtube"`
	Role                 int                           `json:"role" bson:"role"`
	InviteCode           string                        `json:"invite_code" bson:"invite_code"`
	UserInviteCode       string                        `json:"user_invite_code" bson:"user_invite_code"`
	ReferenceCode        string                        `json:"reference_code" bson:"reference_code"`
	AgentCode            string                        `json:"agent_code" bson:"agent_code"`
	FBLink               string                        `json:"fb_link" bson:"fb_link"`
	Email                string                        `json:"email" bson:"email"`
	Status               constants.CommonStatus        `json:"status,omitempty" bson:"status,omitempty"`
	Permissions          []constants.PermissionCommand `json:"permissions,omitempty" bson:"permissions,omitempty"`
	CreatedAt            interface{}                   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	SyncDate             *time.Time                    `json:"sync_date,omitempty" bson:"sync_date,omitempty"`
	EnableAPI8           bool                          `json:"enable_api8" bson:"enable_api8"`
	IsDeCheckpoint       bool                          `json:"is_decheckpoint" bson:"is_decheckpoint"`
	QuotaCloneReg        int                           `json:"quota_clone_reg" bson:"quota_clone_reg"` // User có đăng ký clone reg hay không ? Max clone regs / day
	//
	UpdatedIP       string              `json:"updated_ip,omitempty" bson:"updated_ip,omitempty"`
	UpdatedEmployee *primitive.ObjectID `json:"updated_employee,omitempty" bson:"updated_employee,omitempty"`
	UpdatedUser     *primitive.ObjectID `json:"updated_user,omitempty" bson:"updated_user,omitempty"`
	UpdatedSource   string              `json:"updated_source,omitempty" bson:"updated_source,omitempty"`
	UpdatedDate     *time.Time          `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	//
	CreatedIP       string              `json:"created_ip,omitempty" bson:"created_ip,omitempty"`
	CreatedEmployee *primitive.ObjectID `json:"created_employee,omitempty" bson:"created_employee,omitempty"`
	CreatedUser     *primitive.ObjectID `json:"created_user,omitempty" bson:"created_user,omitempty"`
	CreatedSource   string              `json:"created_source,omitempty" bson:"created_source,omitempty"`
	CreatedDate     *time.Time          `json:"created_date,omitempty" bson:"created_date,omitempty"`
	//
	LastChangePasswordDate *time.Time `json:"-" bson:"last_change_password_date,omitempty"`
}

// IsExists ..
func (m UserInfo) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}

// ActionProfileDefault ..
type ActionProfileDefault struct {
	Key             constants.AppName  `json:"key,omitempty" bson:"key"`
	ActionProfileID primitive.ObjectID `json:"action_profile_id,omitempty" bson:"action_profile_id"`
}

// AllUserReq : get all user
type AllUserReq struct {
	Status   []constants.CommonStatus
	Username string
	Offset   primitive.ObjectID
	Limit    int
}
