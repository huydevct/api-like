package constants

type (
	// CommonStatus : trạng thái chung
	CommonStatus string

	// PackageType : Loại package
	PackageType string

	// TransStatus : Status của transaction
	TransStatus string

	// CommentStatus : Status của comment
	CommentStatus string

	// ServiceStatus : Status của service
	ServiceStatus string

	// ServiceLogStatus : Status của service log
	ServiceLogStatus string

	// GiftStatus : Status của gift
	GiftStatus string

	// AliveStatus : status của clone
	AliveStatus string

	// DebugModeStatus : status của debug mode
	DebugModeStatus string

	// WalletLogType : Loại action wallet
	WalletLogType string

	// CloneActivityStatus : Status của Clone activity
	CloneActivityStatus string

	// HotMailStatus : Status hotmail
	HotMailStatus string

	// TemplateAction : template action
	TemplateAction int

	// DoResultStatus : trạng thái ghi nhận từ mobile
	DoResultStatus string
)

// giá trị mặc định khi đăng ký user
const (
	Balance = 12000
)

// Working status
const (
	Active   CommonStatus = "Active"   // Đang làm việc, hoạt đông
	Pause    CommonStatus = "Pause"    // Đã nghỉ, ngưng hoạt động
	Approved CommonStatus = "Approved" // Đã đuợc duyệt
	Delete   CommonStatus = "Delete"   // Đã xóa
)

// Status clone
const (
	CloneFree               AliveStatus = "free" // Clone Reg, chưa binding vào token
	CloneLive               AliveStatus = "live"
	CloneGetting            AliveStatus = "getting"
	CloneChecking           AliveStatus = "checking"
	CloneCheckpoint         AliveStatus = "checkpoint"
	CloneStored             AliveStatus = "stored"
	CloneDelete             AliveStatus = "delete"
	CloneCheckpointRejected AliveStatus = "checkpoint_rejected"
	CloneWaitingReview      AliveStatus = "waiting_review"
)

// Status hotmail
const (
	HotMailLive HotMailStatus = "Live"
	HotMailUsed HotMailStatus = "Used"
)

// Template action
const (
	TemplateActionAdmin TemplateAction = 1 // aciton dành cho admin
	TemplateActionUser  TemplateAction = 2 // action dành cho user
)

// Status debug mode
const (
	DebugModeProduction DebugModeStatus = "production"
	DebugModeTest       DebugModeStatus = "test"
	DebugModeFullTest   DebugModeStatus = "full_test"
)

// Clone activity
const (
	CloneActivityCreated  CloneActivityStatus = "Created"
	CloneActivityFinished CloneActivityStatus = "Finished"
)

// Enum package status
const (
	PackageDeposit PackageType = "deposit"
	PackagePaypal  PackageType = "paypal"
)

// Transaction status
const (
	TransPending TransStatus = "Pending"
	TransActive  TransStatus = "Active"
)

// Comment status
const (
	CommentPending CommentStatus = "Pending"
	CommentActive  CommentStatus = "Active"
	CommentCancel  CommentStatus = "Cancel"
	CommentDelete  CommentStatus = "Delete"
)

// Service status
const (
	ServicePending         ServiceStatus = "Pending"          // Đang chờ xác nhận mua
	ServiceActive          ServiceStatus = "Active"           // Đang thực hiện
	ServiceWaiting         ServiceStatus = "Waiting"          // Chờ hủy
	ServicePause           ServiceStatus = "Pause"            // Hủy
	ServiceSuccess         ServiceStatus = "Success"          // Đã hoàn thành
	ServiceReport          ServiceStatus = "Report"           // Gói bị autofarm report
	ServicePendingWarranty ServiceStatus = "Pending_Warranty" // Đang chờ xác nhận bảo hành
)

// Servicelog status
const (
	ServicelogCreated  ServiceLogStatus = "Created"  // Vừa tạo, chưa binding
	ServicelogUpdating ServiceLogStatus = "Updating" // Đang chờ binding, re-binding
	ServicelogBinding  ServiceLogStatus = "Binding"  // Đã binding với uid, chờ chạy, chuyển thành getting
	ServicelogGetting  ServiceLogStatus = "Getting"  // Đã trả về cho mobile, chờ chuyển thành finished
	ServicelogFinished ServiceLogStatus = "Finished" // Đã thực hiện xong
)

// Gift status
const (
	GiftActive  GiftStatus = "Active"  // Khả dụng
	GiftWaiting GiftStatus = "Waiting" // Không khả dụng
	GiftUsed    GiftStatus = "Used"    // Đã sử dụng
	GiftDelete  GiftStatus = "Delete"  // Đã xóa
)

// WalletLogType
const (
	WallletRecharge WalletLogType = "recharge" // nạp tiền
	WallletRefund   WalletLogType = "refund"   // hoàn tiền
	WallletGift     WalletLogType = "gift"     // khuyến mãi
	WallletPurchase WalletLogType = "purchase" // mua
)

// DoResultStatus
const (
	DoResultSuccess DoResultStatus = "Success"
	DoResultReport  DoResultStatus = "Report"
)
