package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// User middleware
var (
	User = user{}
)

type user struct{}

// ValidateToken : Kiểm tra token có hợp lệ hay không ?
func (user) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		token := c.Request().Header.Get("Token")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is required!")
		}
		// validate token
		userType := "user"
		userID, userToken, err := ValidateToken(token, userType)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
		}

		// Set cho handler
		c.Set("token", token)              // string
		c.Set("user_id", userID)           // int
		c.Set("user_id_str", userID.Hex()) // string
		c.Set("user_token", userToken)     // string

		return next(c)
	}
}
