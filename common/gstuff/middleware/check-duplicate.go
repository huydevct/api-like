package middleware

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"app/common/gstuff/handler"

	"github.com/labstack/echo/v4"
)

func checkDuplicate(next echo.HandlerFunc, duration time.Duration) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// get request body
		req := c.Request()
		bodyRequest := []byte{}
		if req.Body != nil { // Read
			bodyRequest, _ = ioutil.ReadAll(req.Body)
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRequest)) // Reset

		// hash req body
		hashBody := sha256.Sum256([]byte(bodyRequest))

		keyRedis := fmt.Sprintf("middleware:checkDup:%s/%x", req.RequestURI, hashBody)
		redisClient := cfg.Redis.Get("core").GetClient()
		_, errReq := redisClient.Get(keyRedis).Result()

		if errReq == nil {

			return c.JSON(http.StatusTooManyRequests,
				&handler.ResponseContent{
					Code:    http.StatusTooManyRequests,
					Message: "Too many requests",
				})
		}

		// set new keyRedis
		redisClient.Set(keyRedis, 1, duration)

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}

// CheckDuplicate15Second ..
func CheckDuplicate15Second(next echo.HandlerFunc) echo.HandlerFunc {
	return checkDuplicate(next, time.Second*15)
}

// CheckDuplicate1Min ..
func CheckDuplicate1Min(next echo.HandlerFunc) echo.HandlerFunc {
	return checkDuplicate(next, time.Minute*1)
}

// CheckDuplicate5Min ..
func CheckDuplicate5Min(next echo.HandlerFunc) echo.HandlerFunc {
	return checkDuplicate(next, time.Minute*5)
}

// CheckDuplicate30Day ..
func CheckDuplicate30Day(next echo.HandlerFunc) echo.HandlerFunc {
	return checkDuplicate(next, time.Hour*720)
}
