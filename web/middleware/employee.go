package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Employee middleware
var (
	Employee = employee{}
)

type employee struct{}

// ValidateToken : Kiểm tra token có hợp lệ hay không ?
func (employee) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		token := c.Request().Header.Get("Token")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is required!")
		}
		// validate token
		userType := "employee"
		employeeID, _, err := ValidateToken(token, userType)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
		}

		// Set cho handler
		c.Set("token", token)                      // string
		c.Set("employee_id", employeeID)           // objectID
		c.Set("employee_id_str", employeeID.Hex()) // string

		return next(c)
	}
}

// ValidateMobileSecretkey : Kiểm tra mobile secretkey ?
func (employee) ValidateMobileSecretkey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		secretkey := c.Request().Header.Get("mobile-secret-key")
		if secretkey == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Secret key is required!")
		}
		// validate secretkey

		if secretkey != cfg.Other.Get("mobile-secret-key") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Secret key is invalid!!")
		}

		// Set cho handler
		c.Set("mobile-secret-key", secretkey) // string

		return next(c)
	}
}
