package route

import (
	"github.com/labstack/echo/v4"

	groute "app/common/gstuff/route"
	"app/web/handler"
	"app/web/middleware"
)

// AutofarmerNetAPI : bộ api của user dùng cho web, mobile
// xác thực token
func AutofarmerNetAPI(e *echo.Echo) {

	base := groute.BaseRoute(e)

	v1API := base.Group("/v1")
	{

		v1API.POST("/config", handler.NewConfig().GetConfig)

		v1API.POST("/do-actions", handler.NewActionsHandler().GetDoActions)
		v1API.POST("/update-service", handler.NewServiceHandler().Update)

		apiEmployeeGr := v1API.Group("/employees")
		{
			apiEmployeeGr.POST("/register", handler.NewEmployeeHandler().Register)
		}

		// package group
		publicPackGr := v1API.Group("/packs")
		{
			publicPackGr.POST("/all", handler.NewPackHandler().All)
			publicPackGr.GET("/detail", handler.NewPackHandler().Detail)
		}
		// setting group
		publicSettingGr := v1API.Group("/settings")
		{
			publicSettingGr.GET("/detail", handler.NewSettingHandler().Detail)
		}

		v1API.POST("/register", handler.NewUserHandler().Register)
		v1API.POST("/login", handler.NewUserHandler().Login)
		v1API.POST("/send-otp-reset-password", handler.NewUserHandler().SendOTPResetPassword)
		v1API.POST("/reset-password", handler.NewUserHandler().ResetPassword)
		v1API.POST("/change-password", handler.NewUserHandler().ChangePassword)

		// qrcode group
		qrCodeGr := v1API.Group("/qrcode")
		{
			qrCodeGr.POST("/gen", handler.NewQrCodeHandler().Gen, middleware.User.ValidateToken)
			qrCodeGr.POST("/parse", handler.NewQrCodeHandler().ParseQRCode)
			qrCodeGr.POST("/parse-login", handler.NewQrCodeHandler().ParseQRCodeLogin)
		}
		userGr := v1API.Group("/user")
		{
			userGr.POST("/check-exists", handler.NewUserHandler().CheckUserExists)
		}

		// log group
		logGr := v1API.Group("/logs")
		{
			logGr.POST("/write", handler.NewLogHandler().Write)
		}
	}

	v1APIAuthen := base.Group("/v1", middleware.User.ValidateToken)
	{
		v1APIAuthen.POST("/logout", handler.NewUserHandler().Logout)
		v1APIAuthen.POST("/update-info", handler.NewUserHandler().UpdateInfo)
		v1APIAuthen.GET("/detail", handler.NewUserHandler().DetailByToken)
		v1APIAuthen.POST("/change-password", handler.NewUserHandler().ChangePassword)

		// transaction
		publicTransactionGr := v1APIAuthen.Group("/transactions")
		{
			publicTransactionGr.POST("/create", handler.NewTransactionHandler().Create)
			publicTransactionGr.GET("/detail", handler.NewTransactionHandler().Detail)
			publicTransactionGr.POST("/all", handler.NewTransactionHandler().All)
		}

		// gift code
		publicGiftGr := v1APIAuthen.Group("/gift")
		{
			publicGiftGr.POST("/apply", handler.NewGiftCodeHandler().Apply)
		}

		// wallet logs
		publicWalletLogGr := v1APIAuthen.Group("/wallet-logs")
		{
			publicWalletLogGr.POST("/all", handler.NewWalletLogHandler().All)
		}

		// action profile
		publicActionGr := v1APIAuthen.Group("/action-profile")
		{
			publicActionGr.POST("/create", handler.NewActionProfileHandler().Create)
			publicActionGr.POST("/all", handler.NewActionProfileHandler().All)
			publicActionGr.POST("/update", handler.NewActionProfileHandler().Update)
			publicActionGr.GET("/detail", handler.NewActionProfileHandler().Detail)
			publicActionGr.POST("/delete", handler.NewActionProfileHandler().Delete)
			publicActionGr.POST("/all-by-clone", handler.NewActionProfileHandler().AllByClone)
		}

		// page
		publicPageGr := v1APIAuthen.Group("/pages")
		{
			publicPageGr.POST("/create", handler.NewPageHandler().Create)
			publicPageGr.POST("/all", handler.NewPageHandler().All)
			publicPageGr.GET("/detail", handler.NewPageHandler().Detail)
			publicPageGr.POST("/delete", handler.NewPageHandler().Delete)
		}

		// clone
		publicCloneGr := v1APIAuthen.Group("/clones")
		{
			publicCloneGr.GET("/all", handler.NewCloneHandler().AllClone)
			publicCloneGr.POST("/create", handler.NewCloneHandler().Create)
			publicCloneGr.POST("/set-action", handler.NewCloneHandler().SetAction)
			publicCloneGr.POST("/delete", handler.NewCloneHandler().Delete)
			publicCloneGr.POST("/reset", handler.NewCloneHandler().Reset)
			publicCloneGr.POST("/detail", handler.NewCloneHandler().Detail)
			publicCloneGr.POST("/search", handler.NewCloneHandler().Search)
		}

		// group
		publicGroupGr := v1APIAuthen.Group("/groups")
		{
			publicGroupGr.POST("/create", handler.NewGroupHandler().Create)
			publicGroupGr.POST("/all", handler.NewGroupHandler().All)
			publicGroupGr.GET("/detail", handler.NewGroupHandler().Detail)
			publicGroupGr.POST("/delete", handler.NewGroupHandler().Delete)
		}

		// group uid
		publicGroupUIDGr := v1APIAuthen.Group("/group-uid")
		{
			publicGroupUIDGr.POST("/all", handler.NewGroupUIDHandler().All)
			publicGroupUIDGr.POST("/delete", handler.NewGroupUIDHandler().Delete)
		}

		// friend
		publicFriendGr := v1APIAuthen.Group("/friends")
		{
			publicFriendGr.POST("/create", handler.NewFriendHandler().Create)
			publicFriendGr.POST("/all", handler.NewFriendHandler().All)
			publicFriendGr.GET("/detail", handler.NewFriendHandler().Detail)
			publicFriendGr.POST("/delete", handler.NewFriendHandler().Delete)
		}

		// friend uid
		publicFriendUIDGr := v1APIAuthen.Group("/friend-uid")
		{
			publicFriendUIDGr.POST("/all", handler.NewFriendUIDHandler().All)
			publicFriendUIDGr.POST("/delete", handler.NewFriendUIDHandler().Delete)
		}

		// device
		publicDeviceGr := v1APIAuthen.Group("/devices")
		{
			publicDeviceGr.POST("/delete", handler.NewDeviceHandler().Delete)
			publicDeviceGr.POST("/all", handler.NewDeviceHandler().All)
			publicDeviceGr.GET("/detail", handler.NewDeviceHandler().Detail)
			publicDeviceGr.POST("/search", handler.NewDeviceHandler().Search)
			publicDeviceGr.POST("/total-all-clone", handler.NewDeviceHandler().TotalAllClone)
		}
	}

	v1AdminAPI := base.Group("/v1/admin/")
	{
		v1AdminAPI.POST("/login", handler.NewEmployeeHandler().Login)

		// log group
		logGr := v1AdminAPI.Group("/logs")
		{
			logGr.POST("/write", handler.NewLogHandler().Write)
		}
	}

	v1APIAdminAuthen := base.Group("/v1/admin/", middleware.Employee.ValidateToken)
	{
		v1APIAdminAuthen.POST("/logout", handler.NewEmployeeHandler().Logout)

		// users group
		publicUserGr := v1APIAdminAuthen.Group("/users")
		{
			publicUserGr.POST("/all", handler.NewEmployeeHandler().AllUser)
			publicUserGr.GET("/detail-by-phone", handler.NewEmployeeHandler().DetailUserByPhone)
			publicUserGr.GET("/detail", handler.NewEmployeeHandler().DetailUserByID)
			publicUserGr.POST("/update", handler.NewEmployeeHandler().UpdateUser)
			publicUserGr.POST("/delete", handler.NewEmployeeHandler().DeleteUser)
			publicUserGr.POST("/set-token", handler.NewEmployeeHandler().SetUserToken)
			publicUserGr.POST("/get-token-new-by-old", handler.NewEmployeeHandler().GetTokenNewByOld)
		}

		publicActionGr := v1APIAdminAuthen.Group("/action-profiles")
		{
			publicActionGr.POST("/create", handler.NewActionProfileTemplateHandler().Create)
			publicActionGr.POST("/all", handler.NewActionProfileTemplateHandler().All)
			publicActionGr.POST("/update", handler.NewActionProfileTemplateHandler().Update)
			publicActionGr.GET("/detail", handler.NewActionProfileTemplateHandler().Detail)
			publicActionGr.POST("/delete", handler.NewActionProfileTemplateHandler().Delete)
		}

		// setting group
		publicSettingGr := v1APIAdminAuthen.Group("/settings")
		{
			publicSettingGr.POST("/create", handler.NewSettingHandler().CreateUpdate)
			publicSettingGr.GET("/detail", handler.NewSettingHandler().Detail)
		}

		// setting price group
		publicSettingPriceGr := v1APIAdminAuthen.Group("/setting-prices")
		{
			publicSettingPriceGr.POST("/create", handler.NewSettingPriceHandler().Create)
			publicSettingPriceGr.POST("/update", handler.NewSettingPriceHandler().Update)
			publicSettingPriceGr.GET("/detail", handler.NewSettingPriceHandler().Detail)
			publicSettingPriceGr.POST("/all", handler.NewSettingPriceHandler().All)
		}

		// package group
		publicPackGr := v1APIAdminAuthen.Group("/packs")
		{
			publicPackGr.POST("/create", handler.NewPackHandler().Create)
			publicPackGr.POST("/update", handler.NewPackHandler().Update)
			publicPackGr.POST("/delete", handler.NewPackHandler().Delete)
			publicPackGr.POST("/all", handler.NewPackHandler().All)
			publicPackGr.GET("/detail", handler.NewPackHandler().Detail)
		}

		// transaction group
		publicTransactionGr := v1APIAdminAuthen.Group("/transactions")
		{
			publicTransactionGr.GET("/detail", handler.NewTransactionHandler().Detail)
			publicTransactionGr.POST("/all", handler.NewTransactionHandler().All)
			publicTransactionGr.POST("/active", handler.NewTransactionHandler().Active)
		}

		// gift code
		publicGiftGr := v1APIAdminAuthen.Group("/gift")
		{
			publicGiftGr.POST("/gen", handler.NewGiftCodeHandler().Gen)
			publicGiftGr.POST("/apply", handler.NewGiftCodeHandler().Apply)
			publicGiftGr.POST("/all", handler.NewGiftCodeHandler().All)
			publicGiftGr.POST("/all-history", handler.NewGiftCodeHandler().AllHistory)
			publicGiftGr.POST("/update-status", handler.NewGiftCodeHandler().UpdateStatus)
			publicGiftGr.POST("/delete", handler.NewGiftCodeHandler().Delete)
		}

		// wallet logs
		publicWalletLogGr := v1APIAdminAuthen.Group("/wallet-logs")
		{
			publicWalletLogGr.POST("/all", handler.NewWalletLogHandler().All)
		}

		//device
		publicDevice := v1APIAdminAuthen.Group("/device")
		{
			publicDevice.POST("/all", handler.NewDeviceHandler().AllDevice)
			publicDevice.POST("/detail", handler.NewDeviceHandler().Detail)
			publicDevice.POST("/all-clone", handler.NewDeviceHandler().TotalAllClone)
		}

	}
}
