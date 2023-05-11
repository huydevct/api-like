package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Mobile middleware
var (
	Mobile = mobile{}
)

type mobile struct{}

// ValidateMd5Token : md5 token của user xem có hợp lệ hay ko ?
func (mobile) ValidateMd5Token(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		secretkey := c.Request().Header.Get("mobile-secret-key")
		if secretkey == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Secret key is required!")
		}
		// validate secretkey
		type myRequest struct {
			Token string `json:"token" query:"token"`
		}

		request := myRequest{}
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		}
		json.Unmarshal(bodyBytes, &request)
		// write back to request body
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// So sánh md5(token) == secretkey hay không ?
		data := []byte(request.Token)
		md5Token := fmt.Sprintf("%x", md5.Sum(data))

		if md5Token != secretkey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Secret key is invalid!!")
		}

		return next(c)
	}
}
