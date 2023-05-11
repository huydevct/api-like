package route

import (
	"github.com/labstack/echo/v4"

	groute "app/common/gstuff/route"
	"app/web/handler"
	"app/web/middleware"
)

// AdminAutofarmerNetAPI : bộ api của admin dùng cho web, mobile
// xác thực employee token
func AdminAutofarmerNetAPI(e *echo.Echo) {

	base := groute.BaseRoute(e)

	v1API := base.Group("/v1")
	{
		v1API.POST("/login", handler.NewEmployeeHandler().Login)

		// log group
		logGr := v1API.Group("/logs")
		{
			logGr.POST("/write", handler.NewLogHandler().Write)
		}
	}

	v1APIAuthen := base.Group("/v1", middleware.Employee.ValidateToken)
	{
		v1APIAuthen.POST("/logout", handler.NewEmployeeHandler().Logout)

		// users group
		publicUserGr := v1APIAuthen.Group("/users")
		{
			publicUserGr.POST("/all", handler.NewEmployeeHandler().AllUser)
			publicUserGr.GET("/detail-by-phone", handler.NewEmployeeHandler().DetailUserByPhone)
			publicUserGr.GET("/detail", handler.NewEmployeeHandler().DetailUserByID)
			publicUserGr.POST("/update", handler.NewEmployeeHandler().UpdateUser)
			publicUserGr.POST("/delete", handler.NewEmployeeHandler().DeleteUser)
			publicUserGr.POST("/set-token", handler.NewEmployeeHandler().SetUserToken)
			publicUserGr.POST("/get-token-new-by-old", handler.NewEmployeeHandler().GetTokenNewByOld)
		}

		publicActionGr := v1APIAuthen.Group("/action-profiles")
		{
			publicActionGr.POST("/create", handler.NewActionProfileTemplateHandler().Create)
			publicActionGr.POST("/all", handler.NewActionProfileTemplateHandler().All)
			publicActionGr.POST("/update", handler.NewActionProfileTemplateHandler().Update)
			publicActionGr.GET("/detail", handler.NewActionProfileTemplateHandler().Detail)
			publicActionGr.POST("/delete", handler.NewActionProfileTemplateHandler().Delete)
		}

		// setting group
		publicSettingGr := v1APIAuthen.Group("/settings")
		{
			publicSettingGr.POST("/create", handler.NewSettingHandler().CreateUpdate)
			publicSettingGr.GET("/detail", handler.NewSettingHandler().Detail)
		}

		// setting price group
		publicSettingPriceGr := v1APIAuthen.Group("/setting-prices")
		{
			publicSettingPriceGr.POST("/create", handler.NewSettingPriceHandler().Create)
			publicSettingPriceGr.POST("/update", handler.NewSettingPriceHandler().Update)
			publicSettingPriceGr.GET("/detail", handler.NewSettingPriceHandler().Detail)
			publicSettingPriceGr.POST("/all", handler.NewSettingPriceHandler().All)
		}

		// package group
		publicPackGr := v1APIAuthen.Group("/packs")
		{
			publicPackGr.POST("/create", handler.NewPackHandler().Create)
			publicPackGr.POST("/update", handler.NewPackHandler().Update)
			publicPackGr.POST("/delete", handler.NewPackHandler().Delete)
			publicPackGr.POST("/all", handler.NewPackHandler().All)
			publicPackGr.GET("/detail", handler.NewPackHandler().Detail)
		}

		// transaction group
		publicTransactionGr := v1APIAuthen.Group("/transactions")
		{
			publicTransactionGr.GET("/detail", handler.NewTransactionHandler().Detail)
			publicTransactionGr.POST("/all", handler.NewTransactionHandler().All)
			publicTransactionGr.POST("/active", handler.NewTransactionHandler().Active)
		}

		// gift code
		publicGiftGr := v1APIAuthen.Group("/gift")
		{
			publicGiftGr.POST("/gen", handler.NewGiftCodeHandler().Gen)
			publicGiftGr.POST("/apply", handler.NewGiftCodeHandler().Apply)
			publicGiftGr.POST("/all", handler.NewGiftCodeHandler().All)
			publicGiftGr.POST("/all-history", handler.NewGiftCodeHandler().AllHistory)
			publicGiftGr.POST("/update-status", handler.NewGiftCodeHandler().UpdateStatus)
			publicGiftGr.POST("/delete", handler.NewGiftCodeHandler().Delete)
		}

		// wallet logs
		publicWalletLogGr := v1APIAuthen.Group("/wallet-logs")
		{
			publicWalletLogGr.POST("/all", handler.NewWalletLogHandler().All)
		}

		//device
		publicDevice := v1APIAuthen.Group("/device")
		{
			publicDevice.POST("/all", handler.NewDeviceHandler().AllDevice)
			publicDevice.POST("/detail", handler.NewDeviceHandler().Detail)
			publicDevice.POST("/all-clone", handler.NewDeviceHandler().TotalAllClone)
		}

	}
}
