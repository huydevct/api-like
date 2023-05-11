package constants

//
const (
	// User
	NumberLoginFailKey            = "Autofarmer:NumberLoginFail:%v"
	NumberChangePasswordFailKey   = "Autofarmer:NumberChangePasswordFail:%v"
	TokenLoginKey                 = "Autofarmer:TokenLogin:%v"
	NumberSendOTPResetPasswordKey = "Autofarmer:NumberSendOTPResetPassword:%v"
	OTPKey                        = "Autofarmer:OTP:%v:%v"            // Autofarmer:OTP:{{Type}}:{{Phone}}
	NumberCloneRegKey             = "Autofarmer:NumberCloneReg:%v:%v" // Autofarmer:NumberCloneReg:{{Token}}:{{Day}}

	// ActionProfile
	ActionProfileKey = "Autofarmer:ActionProfile:%v" // Autofarmer:ActionProfile:{{ActionProfileID}}

	// Page
	PageKey = "Autofarmer:Page:%v" // Autofarmer:Page:{{PageID}}

	// User
	UserTokenKey = "Autofarmer:User:TokenNew:%v" //Autofarmer:User:Token:{{token}}

	// SettingPrices
	SettingPricesKey = "Autofarmer:SettingPrices:%v" // Autofarmer:SettingPrices:{{appname}}

	// Setting
	SettingKey = "Autofarmer:Setting"

	// JSON
	JSONJasmineKey = "Autofarmer:JsonJasmine"

	// Clone
	CloneKey = "Autofarmer:Clone:%v" // Autofarmer:Clone:{{cloneID}}

	// Device
	DeviceKey = "Autofarmer:Device:%v" // Autofarmer:Device:{{DeviceID}}

	// QRLogin
	QRLoginKey = "Autofarmer:QRLogin:%v" // Autofarmer:QRLogin:{{tokenQR}}

	// MigrateBindingCloneDeviceID
	MigrateBindingCloneDeviceID = "Autofarmer:Migrate:BindingCloneDeviceID:%v:%v"

	// BindingCloneNewDeviceID
	MigrateBindingCloneNewDeviceID = "Autofarmer:Migrate:BindingCloneNewDeviceID:%v:%v"

	// Jasmine Api5
	JasmineApi5Key = "Api5:Jasmine"
)
