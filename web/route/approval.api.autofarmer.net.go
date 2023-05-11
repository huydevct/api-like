package route

import (
	"github.com/labstack/echo/v4"

	groute "app/common/gstuff/route"
	"app/web/handler"
	"app/web/middleware"
)

// ApprovalAutofarmerNetAPI : bộ api của admin, cần cắm cert để authen
func ApprovalAutofarmerNetAPI(e *echo.Echo) {

	base := groute.BaseRoute(e)

	v1API := base.Group("/v1")
	{
		apiEmployeeGr := v1API.Group("/employees")
		{
			apiEmployeeGr.POST("/register", handler.NewEmployeeHandler().Register)
		}
		apiTransactionGr := v1API.Group("/transactions")
		{
			apiTransactionGr.POST("/active", handler.NewTransactionHandler().Active, middleware.Employee.ValidateMobileSecretkey)
		}
	}
}
